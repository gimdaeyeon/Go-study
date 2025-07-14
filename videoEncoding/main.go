package main

import (
	"bytes"
	"compress/flate"
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

// 코드 실행
// cat video.rgb24 | go run main.go
// 결과 재생
// ffplay -f rawvideo -pixel_format rgb24 -video_size 384x216 -framerate 25 decoded.rgb24

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

	encoded := make([][]byte, len(frames))
	for i := range frames {
		// 다음으로 각 프레임 사이의 델타를 계산하여 데이터를 단순화 한다.
		// 많은 경우 프레임 사이의 픽셀은 크게 변하지 않는다. 따라서 델타의 대부분은 작다.
		// 이러한 작은 델타를 더 효율적으로 저장할 수 있다.

		// 물론 첫 번째 프레임에는 이전 프레임이 없으므로 전체를 저장한다.
		// 이를 키프레임라고 한다. 실제로 키프레임은 주기적으로 계산되며 메타데이터에 구분되어 있다.
		// 키프레임을 압축할 수도 있지만, 나중에 다루겠다.
		// 인코더에서는 (관례에 따라) 프레임 0을 키프레임으로 지정한다.

		// 나머지 프레임은 이전 프레임을 기준으로 델타를 적용한다.
		// 이를 예측 프레임이라고 하며 P-프레임이라고도 한다.

		if i == 0 {
			encoded[i] = frames[i]
			continue
		}

		delta := make([]byte, len(frames[i]))
		for j := 0; j < len(delta); j++ {
			delta[j] = frames[i][j] - frames[i-1][j]
		}

		// 이제 델타 프레임이 있는데, 출력해 보면 0이 여러 개 포함되어 있다.
		// 이런 0 값들은 압축하기에 매우 적합하므로, 우리는 이를 run length 인코딩으로 압축할것이다.
		// 이는 값이 반복되는 횟수를 저장한 후 값을 저장하는 간단한 알고리즘이다.

		// 예를 들어, 시퀀스 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0은 4, 0, 12, 1, 4, 0으로 저장된다.
		// run length 인코딩은 최신 코덱에서는 더 이상 사용되지 않지만, 좋은 연습이며
		// 압축이라는 목표를 달성하기에 충분하다.

		var rle []byte
		for j := 0; j < len(delta); {
			// 현재 값이 반복되는 횟수를 센다.
			var count byte
			for count = 0; count < 255 && j+int(count) < len(delta) && delta[j+int(count)] == delta[j]; count++ {
			}

			// 개수와 값을 저장한다.
			rle = append(rle, count)
			rle = append(rle, delta[j])

			j += int(count)
		}

		// RLE 프레임을 저장한다.
		encoded[i] = rle
	}

	rleSize := size(encoded)
	log.Printf("RLE size: %d bytes (%0.2f%% original size)", rleSize, 100*float32(rleSize)/float32(rawSize))

	// 원본 영상 크기의 1/4까지 줄였다. 하지만 더 줄일 수도 있다.
	// 가장 긴 run이 대부분 0으로 채워져 있다는 점에 주목해보자
	// 프레임간 델타가 보통 작기 때문이다.

	// 여기서 어떤 압축 알고리즘을 쓰느냐에 대한 선택 여지가 있지만,
	// 예제를 단순하게 유지하기 위해 표준 라이브러리에 들어 있는
	// DEFLATE 알고리즘을 사용해보자
	// (DEFLATE 구현 코드는 이 시연 범위를 넘어가므로 자세히 다루지는 않음)

	var deflated bytes.Buffer
	w, err := flate.NewWriter(&deflated, flate.BestCompression)
	if err != nil {
		log.Fatal(err)
	}

	for i := range frames {
		if i == 0 {
			// 이 프레임이 키프레임이므로, 원본 프레임을 그대로 기록합니다.
			if _, err := w.Write(frames[i]); err != nil {
				log.Fatal(err)
			}
			continue
		}

		delta := make([]byte, len(frames[i]))
		for j := 0; j < len(delta); j++ {
			delta[j] = frames[i][j] - frames[i-1][j]
		}
		if _, err := w.Write(delta); err != nil {
			log.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		log.Fatal(err)
	}

	deflatedSize := deflated.Len()
	log.Printf("DEFLATE size %d bytes (%0.2f%% original size)", deflatedSize, 100*float32(deflatedSize)/float32(rawSize))

	// DEFLATE단계는 실행하는데 시간이 오래걸린다.
	// 일반적으로 인코더는 디코더보다 훨씬 느리게 실행되는 경향이 있다.
	// 이는 비디오 코덱뿐만 아니라 대부분의 압축 알고리즘에도 해당한다.
	// 인코더가 데이터를 분석하고 압축 방법을 결정하기 위해 많은 작업을 수행해야하기 때문이다.
	// 반면 디코더는 데이터를 읽고 인코더와 반대되는 작업을 수행하는 단순한 루프이다.

	// 여담이지만, 일반적인 JPEG 압축률이 90% 정도라면
	//  ‘차라리 모든 프레임을 JPEG로 인코딩하면 되지 않을까?’ 하고 생각할 수 있다.
	// 맞는 말이긴 하지만, 우리가 위에서 제시한 알고리즘은 JPEG보다 훨씬 단순하다.

	// 또한, DEFLATE 알고리즘은 데이터의 2차원성을 활용하지 않으므로 효율적이지 않다.
	// 실제 환경에서 비디오 코덱은 여기서 구현한 것보다 훨씬 복잡하다.
	// 코덱은 데이터의 2차원성을 활용하고, 더욱 정교한 알고리즘을 사용하며,
	// 실행되는 하드웨어에 최적화되어있다.
	//  예를 들어, H264 코덱은 많은 최신 GPU하드웨어에 구현되어 있다.

	// 이제 인코딩된 비디오가 있으니, 디코딩하여 어떤 결과가 나오는지 확인해보자

	// 먼저 DEFLATE 스트림을 디코딩한다.
	var inflated bytes.Buffer
	r := flate.NewReader(&deflated)
	if _, err := io.Copy(&inflated, r); err != nil {
		log.Fatal(err)
	}
	if err := r.Close(); err != nil {
		log.Fatal(err)
	}

	// 압축해제된 스트림을 프레임 단위로 나눈다.
	decodedFrames := make([][]byte, 0)
	for {
		frame := make([]byte, width*height*3/2)
		if _, err := io.ReadFull(&inflated, frame); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		decodedFrames = append(decodedFrames, frame)
	}

	// 첫 번째 프레임을 제외한 모든 프레임에 대해 이전 프레임을  델타 프레임에 추가해야한다.
	// 이는 인코더에서 수행한 작업과 반대이다.
	for i := range decodedFrames {
		if i == 0 {
			continue
		}
		for j := 0; j < len(decodedFrames[i]); j++ {
			decodedFrames[i][j] += decodedFrames[i-1][j]
		}
	}
	if err := os.WriteFile("decoded.yuv", bytes.Join(decodedFrames, nil), 0644); err != nil {
		log.Fatal(err)
	}

	// 다음으로 각 YUV 프레임을 RGB로 변환한다.
	for i, frame := range decodedFrames {
		Y := frame[:width*height]
		U := frame[width*height : width*height+width*height/4]
		V := frame[width*height+width*height/4:]

		rgb := make([]byte, 0, width*height*3)
		for j := 0; j < height; j++ {
			for k := 0; k < width; k++ {
				y := float64(Y[j*width+k])
				u := float64(U[(j/2)*(width/2)+(k/2)]) - 128
				v := float64(V[(j/2)*(width/2)+(k/2)]) - 128

				r := clamp(y+1.402*v, 0, 255)
				g := clamp(y-0.344*u-0.714*v, 0, 255)
				b := clamp(y+1.772*u, 0, 255)

				rgb = append(rgb, uint8(r), uint8(g), uint8(b))
			}
		}
		decodedFrames[i] = rgb
	}

	// 마지막으로, 디코딩된 비디오를 파일에 작성한다.
	// 이 비디오는 다음 ffplay로 재생할 수 있다.
	// ffplay -f rawvideo -pixel_format rgb24 -video_size 384x216 -framerate 25 decoded.rgb24
	out, err := os.Create("decoded.rgb24")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	for i := range decodedFrames {
		if _, err := out.Write(decodedFrames[i]); err != nil {
			log.Fatal(err)
		}
	}

}

func size(frames [][]byte) int {
	var size int
	for _, frame := range frames {
		size += len(frame)
	}
	return size
}

func clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
