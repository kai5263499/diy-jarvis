harvard-sentences:
	youtube-dl -o "harvard-sentences.%(ext)s" --write-auto-sub --extract-audio --audio-format mp3 https://youtu.be/A9WgeO9FNzE
	sox harvard-sentences.mp3 -r 16k harvard-sentences.wav remix 1-2
	rm harvard-sentences.mp3
