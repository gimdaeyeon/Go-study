package main

import (
	"flag"
	"io"
	"log"
	"os"
)

// 기본적인 비디오 인코더를 만드는 방법에 대한 내용
// 실제 비디오 인코더는 이보다 훨씬 더 복잡하여 99.9% 이상의 압축률을 달성하지만,
// 이 가이드에서는 간단한 인코더로 90% 압축률을 달성하는 방법을 보여준다.
// 기본적으로 비디오 인코딩은 이미지 인코딩과 매우 유사하지만, 시간적으로 압축할 수 있다.
// 이미지 압축은 종종 인간의 눈이 색상의 작은 변화에 둔감하다는 점을 활용하는데,
// 이 인코더에서도 이 점을 활용한다.

// 또한, 이전 기술을 사용하고 더 많은 수학적 계산이 필요한 최신 기술은 다루지 않는다.
// 이 프로젝트에서는 "최적의" 인코딩 방식에 얽매이지 않고
// 비디오 인코딩의 핵심 개념에 집중하기 위함이다.

// 코드를 실행
//   cat video.rgb24 | go run main.go


func main() {
	var width, height int

	// flag 패키지: 명령줄에서 전달된 옵션(플래그)을 정의하고 파싱해서,
	// 프로그램 안의 변수에 그 값을 할당하도록 돕는 표준 라이브러리
	flag.IntVar(&width, "width", 384, "width of the video")
	flag.IntVar(&height, "height", 216, "height of the video")
	flag.Parse() // Parse() 를 통해서 실제로 cli를 통해 선언한 값이 각 변수에 할당된다.

	frames := make([][]byte, 0 ) // make를 통해 slice생성

	for{
		// stdin에서 원시 비디오 프레임을 읽는다. rgb24형식에서는 각 픽셀(r, g, b)이 1바이트이다.
		// 따라서 프레임의 총 크기는 너비 * 높이 * 3 이다.

		frame := make([]byte, width*height*3)

		// 표준 입력 stdin에서 프레임을 읽는다.
		// io.ReadFull로 정확히 프레임 크기만큼 읽어들여 frame 슬라이스에 채워 넣음
		if _, err := io.ReadFull(os.Stdin, frame); err !=nil{
			break;
		}

		frames = append(frames, frame)
	}

	// 이제 우리는 엄청난 양의 메모리를 사용해서 원시 비디오를 얻었다.

	rawSize := size(frames)
	log.Printf("Raw size: %d bytes",rawSize)


}


func size(frames [][]byte) int{
	var size int
	for _, frame := range frames{
		size += len(frame);
	}
	return size
}


