#!/usr/bin/env python3

import speech_recognition as sr

# obtain path to "english.wav" in the same folder as this script
from os import path

lang='ko-KR'

AUDIO_FILE = path.join(path.dirname(path.realpath(__file__)), "dump","raw.wav")

# use the audio file as the audio source
r = sr.Recognizer()
with sr.AudioFile(AUDIO_FILE) as source:
    audio = r.record(source)  # read the entire audio file

print(r.recognize_google(audio, language=lang).replace(" ", ""))