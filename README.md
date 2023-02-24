# wire-pod-korean

This repo is to support the korean language of wire-pod.
It's very tricky solution using google asr api (e.g. python speechrecognition package).
It was tested under ubuntu 20.04. Raspbian on Raspberry pi will also work but not tested.
Here is how to set up.

1. install speechrecognition package from https://pypi.org/project/SpeechRecognition/
  - "sudo pip3 install SpeechRecognition"
  - "sudo pip3 install pyaudio"
  - python3 was tested but python2 will also work
  - portaudio might be requested using "sudo apt-get install portaudio19-dev"
2. Just follow install instruction as described in wire-pod
  - "sudo STT=vosk ./setup.sh"
  - vosk should be selected as asr engine
3. Set hostname
  - During authentification, vector will find the ubuntu server named "escapepod"
  - "sudo hostnamectl set-hostname escapepod"
  - "sudo systemctl restart avahi-daemon"
  - "sudo systemctl enable avahi-daemon"
4. Install vectorx from "https://github.com/fforchino/vectorx"
 - "git clone https://github.com/fforchino/vectorx"
 - "cd vectorx"
 - "sudo setup.sh"
5. Run wire-pod
  - "sudo chipper/start.sh"
6. Server setting
 - console log will show "Wire-pod is not setup. Use the webserver at port 8080 to set up wire-pod."
 - go to "http://localhost:8080" from chrome browser
 - follow the instruction as wire-pod to authentification
7. Enjoy

Tips:
1. You can change the ASR language from chipper/asr.py
 "lang='ko-KR'" <-> "lang='en-US'"
2. You can find the supported korean command from chipper/intent-data/en-US.json
3. You can update the NLU part from chipper/intent-data/en-US.json
 - add more korean to match internal intent

Disclaimer: I am not expert in golang.


# wire-pod

`wire-pod` is fully-featured server software for the Digital Dream Labs [Vector](https://www.digitaldreamlabs.com/pages/meet-vector) robot. It was created thanks to Digital Dream Labs' [open-sourced code](https://github.com/digital-dream-labs), and is primarily based on [chipper](https://github.com/digital-dream-labs/chipper).

It allows voice commands to work with any Vector 1.0 or 2.0 for no fee, including regular production robots. Not just with OSKR/dev-unlocked robots.

## Documentation

Check out the [wiki](https://github.com/kercre123/wire-pod/wiki) for more information on what wire-pod is, a guide on how to install wire-pod, how to develop for it, and for some helpful tips.

## Credits

- [Digital Dream Labs](https://github.com/digital-dream-labs) for saving Vector and for open sourcing chipper which made this possible
- [dietb](https://github.com/dietb) for rewriting chipper and giving tips
- [fforchino](https://github.com/fforchino) for adding many features such as localization and multilanguage, and for helping out
- [xanathon](https://github.com/xanathon) for the publicity and web interface help
- Anyone who has opened an issue and/or created a pull request for wire-pod
