package main

import (
	"bytes"
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

	frames := make([][]byte, 0) // make를 통해 slice생성

	for {
		// stdin에서 원시 비디오 프레임을 읽는다. rgb24형식에서는 각 픽셀(r, g, b)이 1바이트이다.
		// 따라서 프레임의 총 크기는 너비 * 높이 * 3 이다.

		frame := make([]byte, width*height*3)

		// 표준 입력 stdin에서 프레임을 읽는다.
		// io.ReadFull로 정확히 프레임 크기만큼 읽어들여 frame 슬라이스에 채워 넣음
		if _, err := io.ReadFull(os.Stdin, frame); err != nil {
			break
		}

		frames = append(frames, frame)
	}

	// 이제 우리는 엄청난 양의 메모리를 사용해서 원시 비디오를 얻었다.

	rawSize := size(frames)
	log.Printf("Raw size: %d bytes", rawSize)

	for i, frame := range frames {
		// 먼저, 각 프레임을 yuv420 형식으로 변환한다.
		// 각 픽셀은 RGB24형식으로 다음과 같다.
		// +-----------+-----------+-----------+-----------+
		// |           |           |           |           |
		// | (r, g, b) | (r, g, b) | (r, g, b) | (r, g, b) |
		// |           |           |           |           |
		// +-----------+-----------+-----------+-----------+
		// |           |           |           |           |
		// | (r, g, b) | (r, g, b) | (r, g, b) | (r, g, b) |
		// |           |           |           |           |
		// +-----------+-----------+-----------+-----------+  ...
		// |           |           |           |           |
		// | (r, g, b) | (r, g, b) | (r, g, b) | (r, g, b) |
		// |           |           |           |           |
		// +-----------+-----------+-----------+-----------+
		// |           |           |           |           |
		// | (r, g, b) | (r, g, b) | (r, g, b) | (r, g, b) |
		// |           |           |           |           |
		// +-----------+-----------+-----------+-----------+
		//                        ...
		//
		// YUV420 형식은 다음과 같다.
		//
		// +-----------+-----------+-----------+-----------+
		// |  Y(0, 0)  |  Y(0, 1)  |  Y(0, 2)  |  Y(0, 3)  |
		// |  U(0, 0)  |  U(0, 0)  |  U(0, 1)  |  U(0, 1)  |
		// |  V(0, 0)  |  V(0, 0)  |  V(0, 1)  |  V(0, 1)  |
		// +-----------+-----------+-----------+-----------+
		// |  Y(1, 0)  |  Y(1, 1)  |  Y(1, 2)  |  Y(1, 3)  |
		// |  U(0, 0)  |  U(0, 0)  |  U(0, 1)  |  U(0, 1)  |
		// |  V(0, 0)  |  V(0, 0)  |  V(0, 1)  |  V(0, 1)  |
		// +-----------+-----------+-----------+-----------+  ...
		// |  Y(2, 0)  |  Y(2, 1)  |  Y(2, 2)  |  Y(2, 3)  |
		// |  U(1, 0)  |  U(1, 0)  |  U(1, 1)  |  U(1, 1)  |
		// |  V(1, 0)  |  V(1, 0)  |  V(1, 1)  |  V(1, 1)  |
		// +-----------+-----------+-----------+-----------+
		// |  Y(3, 0)  |  Y(3, 1)  |  Y(3, 2)  |  Y(3, 3)  |
		// |  U(1, 0)  |  U(1, 0)  |  U(1, 1)  |  U(1, 1)  |
		// |  V(1, 0)  |  V(1, 0)  |  V(1, 1)  |  V(1, 1)  |
		// +-----------+-----------+-----------+-----------+

		// 이 형식의 요점은 각 픽셀에 필요한 R, G, B 성분 대신
		// 먼저 다른 공간인 Y(휘도)와 UV(색차)로 변환다는 것이다.
		// Y성분은 픽셀의 밝기이고 UV성분은 픽셀의 색상이다.
		// UV 성분은 인접한 4개의 픽셀에서 공유되므로 4개의 픽셀마다 한 번씩만 저장하면된다.
		// 직관적으로 사람의 눈은 색상보다 밝기에 더 민감하기 때문에
		// 각 픽셀의 밝기를 저장한 다음 각 4개의 픽셀의 색상을 저장할 수 있다.
		// 이렇게 하면 이미지 픽셀의 1/4만 저장하면 되므로 공간을 크게 절약할 수 있다.

		// 추가적으로 YUV형식은 YCbCr이라고도 한다.
		// 사실 완전히 맞는 말은 아니지만, 충분히 비슷하며 색상 공간 선택은 완전히 다른 주제이다.

		// 관례적으로 바이트 슬라이스에서는
		// 왼쪽에서 오른쪽으로 읽은 후 위에서 아래로 저장한다.
		// 즉, i행 j열에 있는 픽셀을 찾으려면 인덱스에 있는바이트를 찾는다.
		// (i * width + j ) * 3

		// 실제로는 이미지가 역순으로 처리되므로 크게 중요하지는 않다.
		// 중요한 것은 일관성을 유지하느 것이다.

		Y := make([]byte, width*height)
		U := make([]float64, width*height)
		V := make([]float64, width*height)

		for j := 0; j < width*height; j++ {
			// 픽셀을 RGB에서 YUV로 변환
			r, g, b := float64(frame[3*j]), float64(frame[3*j+1]), float64(frame[3*j+2])

			// 이 계수는 ITU-R 표준에서 가져온 것이다..
			// https://en.wikipedia.org/wiki/YUV#Y%E2%80%B2UV444_to_RGB888_conversion 참조

			// 실제로 실제 계수는 표준에 따라 달라진다.
			// 예시에서는 크게 중요하지 않다. 중요한 점은
			// YUV로 변환하면 색상 공간을 효율적으로 다운샘플링할 수 있다는 것이다.

			y := +0.299*r + 0.587*g + 0.114*b
			u := -0.169*r - 0.331*g + 0.449*b + 128
			v := 0.499*r - 0.418*g - 0.0813*b + 128

			// YUV값을 바이트 슬라이스에 저장한다.
			// 이 슬라이스들은 다음 단계를 조금 더 쉽게 하기 위해 분리되어 있다.
			Y[j] = uint8(y)
			U[j] = u
			V[j] = v
		}

		// 이제 U와 V의 구성요소를 다운샘플링한다.
		// 이는 U와 V구성 요소를 공유하는 4개의 픽셀을 가져와 평균화하는 과정이다.

		// 다운샘플링된 U와 V구성요소를 이 슬라이스에 저장한다.
		uDownsampled := make([]byte, width*height/4)
		vDownsampled := make([]byte, width*height/4)

		for x := 0; x < height; x += 2 {
			for y := 0; y < width; y += 2 {
				// 이 U와 V구성요소를 공유하는 4개 픽셀의 U 및 V 구성요소의평균을 구한다.
				u := (U[x*width+y] + U[x*width+y+1] + U[(x+1)*width+y] + U[(x+1)*width+y+1]) / 4
				v := (V[x*width+y] + V[x*width+y+1] + V[(x+1)*width+y] + V[(x+1)*width+y+1]) / 4

				// 다운샘플링된 U와 V 구성요소를 바이트 슬라이스에 저장한다.
				uDownsampled[x/2*width/2+y/2] = uint8(u)
				vDownsampled[x/2*width/2+y/2] = uint8(v)
			}
		}

		yuvFrame := make([]byte, len(Y)+len(uDownsampled)+len(vDownsampled))

		// 이제YUV 값을 바이트 슬라이스에 저장해야한다.
		// 데이터 압축률을 높이기 위해 모든 Y값을 먼저 저장하고,
		// 그 다음 모든 U값, 그리고 모든 V 값을 저장한다. 이를 평면 형식이라고 한다.
		// 직관적으로, 인접한 Y, U, V 값은 같은 픽셀에서의 Y, U, V값 자체보다 유사할 가능성이 더 높다.
		// 따라서 구성 요소를 평면 형식으로 저장하면 나중에 더 많은 데이터를 저장 할 수 있다.
		copy(yuvFrame, Y)
		copy(yuvFrame[len(Y):], uDownsampled)
		copy(yuvFrame[len(Y)+len(uDownsampled):], vDownsampled)

		frames[i] = yuvFrame
	}

	// 이제 공간이 절반으로 줄어든 YUV로 인코딩된 비디오가 생겼다.

	yuvSize := size(frames)
	log.Printf("YUV420P size: %d bytes (%0.2f%% original size)", yuvSize, 100*float32(yuvSize)/float32(rawSize))

	// ffplay로 재생할 수 있는 파일에도 쓸 수 있다.
	// ffplay -f rawvideo -pixel_format yuv420p -video_size 384x216 -framerate 25 encoded.yuv
	if err := os.WriteFile("encoded.yuv", bytes.Join(frames, nil), 0644); err != nil {
		log.Fatal(err)
	}

}

func size(frames [][]byte) int {
	var size int
	for _, frame := range frames {
		size += len(frame)
	}
	return size
}
