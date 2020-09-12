# wav_slicer

This utility takes a WAV file as input, breaks it up into smaller chunks (default 10 second), and sends them through the DIY Jarvis processing system.

I use this primarialy for testing the system without yelling at my microphone all the time and also for transcribing the audio from recorded Zoom conference calls which happen to be the same format (mono, 16 bit, 16000Hz) that the DeepSpeech audio processing service requires.

Todo:
* This utility needs to detect and resample wav files into mono 16bit 16000Hz.