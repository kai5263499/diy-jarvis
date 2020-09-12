# mic_capture

This utility takes input from your default pulseaudio source (microphone)  as input, breaks it up into smaller chunks (default 10 second), sends those chunks to the rest of the DIY Jarvis system, and then prints the voice-to-text results.

Todo:
* Analyze the energy level in the sampled audio skip sending the sampled chunk if its below a threshold. An even better approach would be to increase that threshold (up to a point) based on the response from the audio processor