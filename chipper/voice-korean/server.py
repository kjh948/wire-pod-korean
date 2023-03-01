from bottle import route, run, template, request
from chatgpt_wrapper import ChatGPT
from asr import runasr, runasr_raw
from tts import runtts_google, runtts_espeak, runtts_vector_sdk_tts, runtts_vector_sdk_wav
bot = ChatGPT()

@route('/chat', method='POST')
def task():
    params = request.json

    print(request.body)

    if params["command"] == "asr":
        print("asr file name = ", params["file"])
        response = runasr(params["file"])        
    elif params["command"] == "asr_raw":
        print("asr file name = ", params["file"])
        response = runasr_raw(params["file"])        
    elif params["command"] == "chatgpt":
        print("chatgpt input = ", params["text"])
        output = bot.ask(params["text"])
        response = output
    elif params["command"] == "chatgpt_tts_wav":
        print("chatgpt input = ", params["text"])
        output = bot.ask("간단하게 응답해줘 " + params["text"])

        output = output.split('.')[0] + ". " + output.split('.')[1]
        runtts_vector_sdk_wav(output)
        response = output        
    elif params["command"] == "tts":
        print("tts input = ", params["text"])
        runtts_google(params["text"])
        response = None
    elif params["command"] == "tts_sdk":
        print("tts input = ", params["text"])
        runtts_vector_sdk_tts(params["text"])
        response = None
    elif params["command"] == "tts_wav":
        print("tts input = ", params["text"])
        runtts_vector_sdk_wav(params["text"])
        response = None
    return response

run(host='localhost', port=8888)