# mic_capture

This utility takes input from your default pulseaudio source (microphone)  as input, breaks it up into smaller chunks (default 10 second), sends those chunks to an audio processing service via GRPC, and then prints the results.

Todo:
* Analyze the energy level in the sampled audio skip sending the sampled chunk if its below a threshold. An even better approach would be to increase that threshold (up to a point) based on the response from the audio processor