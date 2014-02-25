package rhythmbox

// Debug
func (r *Client) Debug() {
	r.Execute("--debug")
}

// Don't start a new instance of Rhythmbox
func (r *Client) NoStart() {
	r.Execute("--no-start")
}

// Quit Rhythmbox
func (r *Client) Quit() {
	r.Execute("--quit")
}

// Check if Rhythmbox is already running
func (r *Client) CheckRunning() {
	r.Execute("--check-running")
}

// Don't present an existing Rhythmbox window
func (r *Client) NoPresent() {
	r.Execute("--no-present")
}

// Jump to next song
func (r *Client) Next() {
	r.Execute("--next")
}

// Jump to previous song
func (r *Client) Previous() {
	r.Execute("--previous")
}

// Seek in current track
func (r *Client) Seek() {
	r.Execute("--seek")
}

// Resume playback if currently paused
func (r *Client) Play() {
	r.Execute("--play")
}

// Pause playback if currently playing
func (r *Client) Pause() {
	r.Execute("--pause")
}

// Toggle play/pause mode
func (r *Client) PlayPause() {
	r.Execute("--play-pause")
}

// Play a specified URI, importing it if necessary
func (r *Client) PlayUri() {
	r.Execute("--play-uri=URI to play")
}

// Add specified tracks to the play queue
func (r *Client) Enqueue(location string) {
	r.Execute("--enqueue", location)
}

// Empty the play queue before adding new tracks
func (r *Client) ClearQueue() {
	r.Execute("--clear-queue")
}

// Print the title and artist of the playing song
func (r *Client) PrintPlaying() {
	r.Execute("--print-playing")
}

// Print formatted details of the song
func (r *Client) PrintPlayingFormat() {
	r.Execute("--print-playing-format")
}

// Select the source matching the specified URI
func (r *Client) SelectSource() {
	r.Execute("--select-source=Source to select")
}

// Activate the source matching the specified URI
func (r *Client) ActivateSource() {
	r.Execute("--activate-source=Source to activate")
}

// Play from the source matching the specified URI
func (r *Client) PlaySource() {
	r.Execute("--play-source=Source to play from")
}

// Enable repeat playback order
func (r *Client) Repeat() {
	r.Execute("--repeat")
}

// Disable repeat playback order
func (r *Client) NoRepeat() {
	r.Execute("--no-repeat")
}

// Enable shuffle playback order
func (r *Client) Shuffle() {
	r.Execute("--shuffle")
}

// Disable shuffle playback order
func (r *Client) NoShuffle() {
	r.Execute("--no-shuffle")
}

// Set the playback volume
func (r *Client) SetVolume() {
	r.Execute("--set-volume")
}

// Increase the playback volume
func (r *Client) VolumeUp() {
	r.Execute("--volume-up")
}

// Decrease the playback volume
func (r *Client) VolumeDown() {
	r.Execute("--volume-down")
}

// Print the current playback volume
func (r *Client) PrintVolume() {
	r.Execute("--print-volume")
}

// Set the rating of the current song
func (r *Client) SetRating() {
	r.Execute("--set-rating")
}

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
