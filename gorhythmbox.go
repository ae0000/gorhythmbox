package main

import (
	"fmt"
	"os/exec"
)

type RbCommand string

const (
	Rb                           = "rhythmbox-client"                     // The actual client to run commands through
	Debug              RbCommand = "--debug"                              // Debug
	NoStart            RbCommand = "--no-start"                           // Don't start a new instance of Rhythmbox
	Quit               RbCommand = "--quit"                               // Quit Rhythmbox
	CheckRunning       RbCommand = "--check-running"                      // Check if Rhythmbox is already running
	NoPresent          RbCommand = "--no-present"                         // Don't present an existing Rhythmbox window
	Next               RbCommand = "--next"                               // Jump to next song
	Previous           RbCommand = "--previous"                           // Jump to previous song
	Seek               RbCommand = "--seek"                               // Seek in current track
	Play               RbCommand = "--play"                               // Resume playback if currently paused
	Pause              RbCommand = "--pause"                              // Pause playback if currently playing
	PlayPause          RbCommand = "--play-pause"                         // Toggle play/pause mode
	PlayUri            RbCommand = "--play-uri=URI to play"               // Play a specified URI, importing it if necessary
	Enqueue            RbCommand = "--enqueue"                            // Add specified tracks to the play queue
	ClearQueue         RbCommand = "--clear-queue"                        // Empty the play queue before adding new tracks
	PrintPlaying       RbCommand = "--print-playing"                      // Print the title and artist of the playing song
	PrintPlayingFormat RbCommand = "--print-playing-format"               // Print formatted details of the song
	SelectSource       RbCommand = "--select-source=Source to select"     // Select the source matching the specified URI
	ActivateSource     RbCommand = "--activate-source=Source to activate" // Activate the source matching the specified URI
	PlaySource         RbCommand = "--play-source=Source to play from"    // Play from the source matching the specified URI
	Repeat             RbCommand = "--repeat"                             // Enable repeat playback order
	NoRepeat           RbCommand = "--no-repeat"                          // Disable repeat playback order
	Shuffle            RbCommand = "--shuffle"                            // Enable shuffle playback order
	NoShuffle          RbCommand = "--no-shuffle"                         // Disable shuffle playback order
	SetVolume          RbCommand = "--set-volume"                         // Set the playback volume
	VolumeUp           RbCommand = "--volume-up"                          // Increase the playback volume
	VolumeDown         RbCommand = "--volume-down"                        // Decrease the playback volume
	PrintVolume        RbCommand = "--print-volume"                       // Print the current playback volume
	SetRating          RbCommand = "--set-rating"                         // Set the rating of the current song
)

// FORMAT OPTIONS
//    %at    album title
//    %aa    album artist
//    %aA    album artist (lowercase)
//    %as    album artist sortname
//    %aS    album artist sortname (lowercase)
//    %ay    album year
//    %ag    album genre
//    %aG    album genre (lowercase)
//    %an    album disc number
//    %aN    album disc number, zero padded
//    %st    stream title
//    %tn    track number (i.e 8)
//    %tN    track number, zero padded (i.e 08)
//    %tt    track title
//    %ta    track artist
//    %tA    track artist (lowercase)
//    %ts    track artist sortname
//    %tS    track artist sortname (lowercase)
//    %td    track duration
//    %te    track elapsed time
//
// Variables  can  be  combined  using  quotes.  For example "%tn %aa %tt", will
// print the track number followed by the artist and the title of the track.

func CallRb(r ...RbCommand) {
	fmt.Println(r...)
	cmd := exec.Command(Rb, string(r))
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(out))
}

func main() {
	CallRb(PlayPause)
	CallRb(PrintPlaying)
	CallRb()
}
