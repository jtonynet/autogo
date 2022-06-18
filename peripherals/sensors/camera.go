package peripherals

import (
	"fmt"
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/hybridgroup/mjpeg"
	"github.com/jtonynet/autogo/config"
	"gocv.io/x/gocv"
)

var (
	deviceID int
	err      error
	webcam   *gocv.VideoCapture
	stream   *mjpeg.Stream
)

func mjpegCapture() {
	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}

		if img.Empty() {
			continue
		}

		gocv.CvtColor(img, &img, 0)

		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf.GetBytes())
	}
}

func CameraServeStream(cfg config.Camera) {

	deviceID = 0 //  -1 ?
	host := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	webcam, err = gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	width := cfg.Width
	height := float64(width) * 0.75

	webcam.Set(3, float64(width))
	webcam.Set(4, float64(height))
	//webcam.Set(19, 1) // fps

	stream = mjpeg.NewStream()

	go mjpegCapture()

	fmt.Println("Capturing, point your browser to", host)

	http.Handle("/", stream)
	log.Fatal(http.ListenAndServe(host, nil))
}
