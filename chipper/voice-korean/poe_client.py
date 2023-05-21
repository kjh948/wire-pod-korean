import poe
import json

with open('../apiConfig.json','r') as f:
    data = json.load(f)    

try:
    client = poe.Client(data['knowledge']['key'])
    bot = 'chinchilla'
    client.send_chat_break(bot)
except:
    error_response = "가이드에 맞춰 poe 키를 입력하세요"


def ask(message):
    try:
        for chunk in client.send_message(bot, message):
            pass
        print(chunk["text"])    
        return chunk["text"]
    except:
        return error_response

if __name__ == "__main__":
    #Auth
    print("The selected bot is : ", bot)
    #---------------------------------------------------------------------------

    print("Context is now cleared")
    while True:
        message = input("Human : ")
        if message =="!clear":
            client.send_chat_break(bot)
            print("Context is now cleared")
            continue
        if message =="!break":
            break
        reply = ask(message)        
        print(f"{bot} : {reply}")