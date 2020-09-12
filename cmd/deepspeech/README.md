# deepspeech

Connects to an MQTT broker for `RegisterAudioSourceRequest` messages from other clients sent to the control topic `requests`. When these are detected, a dedicated MQTT channel is attempted to obtain raw audio data. After processing raw audio data, the resulting text is sent as a `OutputResponse` control message.

# Todo
* Time-out and reap expired audio sources