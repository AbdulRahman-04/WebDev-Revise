import twilio from "twilio"
import config from "config"


const SID =  config.get("SID")
const TOKEN = config.get("TOKEN")
const PHONE =  config.get("PHONE")

let client = new twilio(SID, TOKEN)

async function sendSMS(smsData) {

    try {

        await client.messages.create({
          body: smsData.body,
          to: smsData.to,
          from: PHONE
        })
        
    } catch (error) {
        console.log(error);
        
    }
    
}

export default sendSMS