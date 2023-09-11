package manifest

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"

	dash "github.com/comcast-viper-packager/libdash/v2"
)

func decodeManifest(b []byte) (*dash.Mpd, error) {
	mpd := new(dash.Mpd)

	err := xml.Unmarshal(b, mpd)
	if err != nil {
		return mpd, err
	}
	return mpd, nil

}

var (
	lastMpdTime     dash.XsdDatetime
	lastRecTime     time.Time
	currentSegments []*dash.S
)

// process Mpd runs as a process and manages the checkif MPD as they are received
func ProcessMpd(b []byte) {

	//var lastMpdTime dash.XsdDatetime
	//var lastRecTime time.Time
	pubDiff := time.Duration(time.Millisecond * 0)
	recDiff := time.Duration(time.Millisecond * 0)
	loc, _ := time.LoadLocation("UTC")
	//timeFormat := "20060102_150405.000"

	log.Println("Processing Manifest...")
	receivedTime := time.Now().In(loc)

	mpd, err := decodeManifest(b)
	if err != nil {
		fmt.Printf("Error decoding manfest: %v", err)
		return
	}

	expandSegmentTimeline(mpd)

	//log.Println(mpd)
	log.Println("Publisth Time: ", mpd.PublishTime)
	log.Println("Last Publish Time: ", lastMpdTime.Time)

	// calc time differences
	pubDiff = mpd.PublishTime.Time.Sub(time.Time(lastMpdTime.Time))
	recDiff = receivedTime.Sub(lastRecTime)

	// Write data to logFile
	data := fmt.Sprintf("New Rc Time: %v, diff: %v\n", receivedTime.In(loc), recDiff)
	log.Printf("%s\n", data)
	data = fmt.Sprintf("New Pub Time: %v, diff: %v", mpd.PublishTime, pubDiff)

	if pubDiff != time.Millisecond*1920 {
		data = data + " -- timing Error\n"
	} else {
		data = data + fmt.Sprintln()
	}

	log.Println(data)

	lastMpdTime = *mpd.PublishTime
	lastRecTime = receivedTime

}

func expandSegmentTimeline(m *dash.Mpd) {
	for _, p := range m.Period {
		for _, a := range p.AdaptationSet {
			log.Printf("Period: %s  ASet: %d Content: %s", p.Id, a.Id, a.ContentType)
			if a.SegmentTemplate != nil && a.SegmentTemplate.SegmentTimeline != nil &&
				a.SegmentTemplate.SegmentTimeline.PublishedSegments != nil {
				currentSegments = a.SegmentTemplate.SegmentTimeline.PublishedSegments
			} else {
				break
			}
			log.Println("number of Segments:", len(currentSegments))
			//Adjust the T values for each segment
			for i := 0; i < len(currentSegments); i++ {
				//log.Printf("I= %d T=%d D=%d", i, currentSegments[i].T, currentSegments[i].D)
				if i > 0 && currentSegments[i].T == 0 {
					currentSegments[i].T = currentSegments[i-1].T + currentSegments[i-1].D
				}
			}
			log.Println("Segments: ", currentSegments)
			currentSegments = nil
		}
	}
}

func removeSegments(m *dash.Mpd) {

}

func addSegments(m *dash.Mpd) {

}
