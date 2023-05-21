import poe

with open('poe_cred.txt','r') as f:
    cred = f.readlines()
client = poe.Client(cred[0])

bot = 'chinchilla'
client.send_chat_break(bot)

def ask(message):
    for chunk in client.send_message(bot, message):
        pass
    print(chunk["text"])    
    return chunk["text"]


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