package main

import (
	"github.com/Philiphil/Ddate-go"
	"time"
	"net/http"
	"fmt"
)

type UpdateFunc func(*Handler, DDate.DTime)
type HttpHandler func(w http.ResponseWriter, r *http.Request)

type Handler struct {
	updateFunc UpdateFunc
	outstr string
	endpoint string
}

// See the TxtUpdate/YmlUpdate for an example of usage
func Formatter(fmtNoHoliday string, fmtHoliday string, currTime DDate.DTime)(string) {
        if(currTime.Holyday == "") {
                // no holyday
                return fmt.Sprintf(fmtNoHoliday, currTime.DayN.String(), currTime.Day,
                        currTime.Season.String(), currTime.Year)
        } else {
                return fmt.Sprintf(fmtHoliday, currTime.DayN.String(), currTime.Day,
                        currTime.Season.String(), currTime.Year, currTime.Holyday)
        }
}

func TxtUpdate(h *Handler, currTime DDate.DTime) {
	var noholiday = "Today is %s, the %dth day of %s in the YOLD %d."
	
	if(currTime.Day % 10 == 1 && currTime.Day != 11) {
		noholiday = "Today is %s, the %st day of %s in the YOLD %d."
	} else if(currTime.Day % 10 == 2 && currTime.Day != 12) {
		noholiday = "Today is %s, the %nd day of %s in the YOLD %d."
	} else if(currTime.Day % 10 == 3 && currTime.Day != 13) {
		noholiday = "Today is %s, the %rd day of %s in the YOLD %d."
	}
	
	var holiday = fmt.Sprintf("%s Celebrate %s", noholiday)

	h.outstr = Formatter(noholiday, holiday, currTime)
}

func JsonUpdate(h *Handler, currTime DDate.DTime) {
	h.outstr = Formatter(`{"dayOfWeek":"%s", "day":%d, "season":"%s", "year":%d}`,
	`{"day":%d, "dayOfWeek":"%s", "season":"%s", "year":%d, "holyday":"%s"}`,
	currTime)
}

func XmlUpdate(h *Handler, currTime DDate.DTime) {
	h.outstr = Formatter("<DayOfWeek>%s</DayOfWeek>\n<Day>%d</Day>\n<Season>%s</Season>\n<Year>%d</Year>",
	"<DayOfWeek>%d</DayOfWeek>\n<Day>%s</Day>\n<Season>%s</Season>\n<Year>%d</Year>\n<Holyday>%s</Holyday>",
	currTime)
}

func YmlUpdate(h *Handler, currTime DDate.DTime) {
	h.outstr = Formatter("dayOfWeek: %s\nday: %d\nseason: %s\nyear: %d\n",
	"day: %d\ndayOfWeek: %s\nseason: %s\nyear: %d\nholyday: %s", currTime)
}

var TxtHandler = Handler{TxtUpdate, "", "/txt"}
var JsonHandler = Handler{JsonUpdate, "", "/json"}
var XmlHandler = Handler{XmlUpdate, "", "/xml"}
var YmlHandler = Handler{YmlUpdate, "", "/yml"}

var handlers = [...]Handler{ TxtHandler, JsonHandler, XmlHandler, YmlHandler}

func update() {
	var dtime = DDate.TimeToDTime(time.Now())
	for i := 0; i < len(handlers); i++ {
		handlers[i].updateFunc(&(handlers[i]), dtime)
		fmt.Printf("[update] %s: %s\n", handlers[i].endpoint, handlers[i].outstr)
	}
}

/*
 * Handler type:
 * update()
 * output()
 */
func main() {
	update()
	for i := 0; i < len(handlers); i++ {
		var handler = handlers[i]
		http.HandleFunc(handler.endpoint, func(w http.ResponseWriter,
		r *http.Request){ fmt.Fprintf(w, handler.outstr) })
	}

	var landing = `<html><body>We support the following endpoints:
<code>/txt</code>, <code>/json</code>, <code>/xml</code>, <code>/yml</code>`
	http.HandleFunc("/", 
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, landing)})
	http.ListenAndServe(":8080", nil)
}
