import nodemailer from "nodemailer"
import config from "config"

const EMAIL = config.get("EMAIL");
const PASS = config.get("PASS")


async function sendEmail(emailData) {

    try {
         
        let transporter = nodemailer.createTransport({
            host: "smtp.gmail.com",
            port: 465,
            secure: true,
            auth: {
                user: EMAIL,
                pass: PASS
            }
        })
         
        let sender = await transporter.sendMail({
            from: EMAIL,
            to :emailData.to,
            subject: emailData.subject,
            html: emailData.html
        })
        
    } catch (error) {
        console.log(error);
        
    }
    
}

export default sendEmail;