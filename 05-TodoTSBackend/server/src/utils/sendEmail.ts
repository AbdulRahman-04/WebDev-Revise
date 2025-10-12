import nodemailer from "nodemailer"
import config from "config"

const USER: string = config.get<string>("USER")
const PASS: string = config.get<string>("PASS")

interface EmailData {
    from :string,
    to: string,
    text?: string,
    subject: string,
    html: string
}

async function sendEmail(emailData:EmailData) {

    try {
 
        let transporter = nodemailer.createTransport({
            host: "smptp.gmail.com",
            port: 465,
            secure: true,
            auth: {
                user: USER,
                pass: PASS
            }
        })

        let sender = await transporter.sendMail({
            from: emailData.from,
            to: emailData.to,
            text: emailData.text,
            subject: emailData.subject,
            html: emailData.html
        })

        
    } catch (error) {
        console.log(error);
        
    }
    
}

export default sendEmail;