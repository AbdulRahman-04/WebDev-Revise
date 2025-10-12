import twilio from "twilio"
import config from  "config"

const SID: string = config.get<string>("SID")
const TOKEN: string = config.get<string>("TOKEN")
const PHONE: string = config.get<string>("PHONE")

interface smsData {
    body: string,
     to: string
}

let client = new twilio.Twilio(SID, TOKEN)

async function sendSMS(SMSData: smsData) {
     try {

        await client.messages.create({
            body: SMSData.body,
            to: SMSData.to,
            from: PHONE
        })


        
     } catch (error) {
        console.log(error);
        
     }
}
