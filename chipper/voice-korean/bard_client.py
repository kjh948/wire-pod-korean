from bardapi import Bard
import json

try:
    with open('bard_cred.txt','r') as f:
        cred = f.readlines()
        bard = Bard(token=cred[0])
except:
    with open('../apiConfig.json','r') as f:
        data = json.load(f)    
        bard = Bard(token=cred[0])



def ask(message):    
    return bard.get_answer(message)['content']


if __name__ == "__main__":
    #Auth
    bot = 'bard'
    print("The selected bot is : ", bot)
    #---------------------------------------------------------------------------

    print("Context is now cleared")
    while True:
        message = input("Human : ")
        reply = ask(message)        
        print(f"{bot} : {reply}")