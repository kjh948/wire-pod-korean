#!/usr/bin/env python3

import speech_recognition as sr
from os import path
import soundfile as sf

# use the audio file as the audio source
r = sr.Recognizer()

def runasr(filename, lang='ko-KR'):
    with sr.AudioFile(filename) as source:
        audio = r.record(source)  # read the entire audio file

    result = r.recognize_google(audio, language=lang).replace(" ", "")
    return result



def runasr_raw(filename, lang='ko-KR'):
    data, samplerate = sf.read(filename, channels=1, samplerate=16000, subtype='PCM_16')
    sf.write(filename+".wav", data, samplerate)
    return runasr(filename+".wav")