package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	max_iter = 1000
)

var (
	virt_x0 = -2.0
	virt_w  = 3.0
	virt_y0 = -1.0
	virt_h  = 2.0
	height  = 900
	width   = int(virt_w * float64(height) / float64(virt_h))

	gtkSpinner     *gtk.Spinner
	gtkImg         *gtk.Image
	gtkProgressbar *gtk.ProgressBar
	goImg          *image.RGBA
	progress       int
)

func mandelbrot(x float64, y float64) bool {
	z := complex(0, 0)
	for iter := 0; iter < max_iter; iter++ {
		z = z*z + complex(x, y)
		if cmplx.Abs(z) > 2 {
			return true
		}
	}
	return false
}

func plot() {
	for y := 0; y != height; y++ {
		for x := 0; x < width; x++ {
			var vx float64 = (float64(x)/float64(width))*virt_w + virt_x0
			var vy float64 = (float64(y)/float64(height))*virt_h + virt_y0
			if mandelbrot(vx, vy) {
				goImg.Set(x, y, color.Black)
			} else {
				goImg.Set(x, y, color.White)
			}
		}
		gtkProgressbar.SetFraction(float64(y) / float64(height))
	}
}

func saveGoImgToGTK() {
	file, err := os.Create("res.png")
	if err != nil {
		log.Fatalf("Error opening result file: %v", err)
	}

	if err := png.Encode(file, goImg); err != nil {
		log.Fatalf("Error writing result file: %v", err)
	}

	gtkImg.SetFromFile("res.png")
}

func main() {
	fmt.Printf("Hello Go ! \n")

	goImg = image.NewRGBA(image.Rect(0, 0, width, height))

	// Create Gtk Application, change appID to your application domain name reversed.
	const appID = "com.github.salicorne.mandelgo"
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatalf("Could not create GTK application: %s", err)
	}
	// Application signals available
	// startup -> sets up the application when it first starts
	// activate -> shows the default first window of the application (like a new document). This corresponds to the application being launched by the desktop environment.
	// open -> opens files and shows them in a new window. This corresponds to someone trying to open a document (or documents) using the application from the file browser, or similar.
	// shutdown ->  performs shutdown tasks
	// Setup Gtk Application callback signals
	application.Connect("activate", func() { onGTKActivate(application) })
	os.Exit(application.Run(os.Args))
}

func GTKGenerateBtnAction() {
	startGTKLoading()
	go func() {
		plot()
		saveGoImgToGTK()
		glib.IdleAdd(stopGTKLoading)
	}()
}

func startGTKLoading() bool {
	gtkSpinner.Start()
	gtkProgressbar.SetVisible(true)
	return false
}

func stopGTKLoading() bool {
	gtkSpinner.Stop()
	gtkProgressbar.SetVisible(false)

	return false
}

// Callback signal from Gtk Application
func onGTKActivate(application *gtk.Application) {
	// Create ApplicationWindow
	appWindow, err := gtk.ApplicationWindowNew(application)
	if err != nil {
		log.Fatalf("Could not create application window: %s", err)
	}

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		panic(err)
	}
	box.SetHomogeneous(false)

	evBox, err := gtk.EventBoxNew()
	if err != nil {
		panic(err)
	}
	evBox.SetAboveChild(true)

	gtkImg, err = gtk.ImageNewFromFile("res.png")
	if err != nil {
		log.Fatalf("Error loading res.png as GTK image: %s", err)
	}
	evBox.Add(gtkImg)
	evBox.Connect("button_press_event", func(w *gtk.EventBox, ev *gdk.Event) bool {
		eventMotion := gdk.EventMotionNewFromEvent(ev)
		x, y := eventMotion.MotionVal()
		log.Printf("Clicked on %v %v\n", x, y)
		return false
	})
	box.PackStart(evBox, true, true, 0)

	// Build the action bar
	actionBar, err := gtk.ActionBarNew()
	if err != nil {
		panic(err)
	}
	box.PackStart(actionBar, false, false, 0)

	gtkSpinner, err = gtk.SpinnerNew()
	if err != nil {
		panic(err)
	}
	actionBar.PackStart(gtkSpinner)

	gtkProgressbar, err = gtk.ProgressBarNew()
	if err != nil {
		panic(err)
	}
	gtkProgressbar.SetVAlign(gtk.ALIGN_CENTER)
	gtkProgressbar.SetSizeRequest(300, 0)
	actionBar.PackStart(gtkProgressbar)

	saveBtn, err := gtk.ButtonNewWithLabel("Save")
	if err != nil {
		panic(err)
	}
	saveBtn.Connect("clicked", func() {
		log.Println("Saving current image on disk")
	})
	actionBar.PackEnd(saveBtn)

	genBtn, err := gtk.ButtonNewWithLabel("Generate")
	if err != nil {
		panic(err)
	}
	genBtn.Connect("clicked", func() {
		log.Println("Generating image")
		GTKGenerateBtnAction()
	})
	actionBar.PackEnd(genBtn)

	appWindow.Add(box)

	// Set ApplicationWindow Properties
	appWindow.SetTitle("Mandelgo dev")
	appWindow.SetDefaultSize(width, height)
	appWindow.ShowAll()

	//plot()
	//saveGoImgToGTK()
	stopGTKLoading()
}
