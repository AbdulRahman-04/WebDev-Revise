import nodemailer from "nodemailer"
import config from "config"

const USER = config.get("USER")
const PASS = config.get("PASS")

async function SendEmail(emailData){
 
    try {

        let transporter = nodemailer.createTransport({
            host: "smtp.gmail.com",
            port:465,
            secure: true,
            auth: {
                user: USER,
                pass : PASS
            }
        })

        let sender = transporter.sendMail({
            from: USER,
            to: emailData.to,
            subject: emailData.subject,
            text: emailData.text,
            html: emailData.html
        })

         console.log(`email sent successfully!🙌`);
        
    } catch (error) {
        console.log(error);
        
    }


}
export default SendEmail