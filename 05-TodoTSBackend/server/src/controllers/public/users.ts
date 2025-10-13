import express, {Request, Response} from "express"
import jwt from "jsonwebtoken"
import bcrypt from "bcrypt"
import sendEmail from "../../utils/sendEmail.js"
// import sendSMS from "../../utils/sendSMS"
import config from "config"
import { userModel } from "../../models/users.js"

const router = express.Router()

const URL: string = config.get<string>("URL")
const USER: string = config.get<string>("USER")
const KEY: string = config.get<string>("JWT_KEY")

router.post("/usersignup", async (req: Request, res: Response) : Promise<void> =>{
    try {
        
        const {userName, email, password, age} = req.body;
        
        // check if all fields are filled or not
        if(!userName || !email || !password || !age){
         res.status(400).json({msg: "fill all the fields "});
         return;
        }

        // duplicate check
        let userExist = await userModel.findOne({email})

        if(userExist){
            res.status(200).json({msg: "user already exists, please sign in"})
            return
        }

        // password hash
        const hashPass = await bcrypt.hash(password, 10);

        // random token for email and phone verify
        let emailToken = Math.random().toString(36).substring(2)
        // let phoneToken = Math.random().toString(36).substring(2)

        // create new obj 

        let newUser = {
            userName,
            email,
            password: hashPass,
            age,
            userVerifyToken: {
                email: emailToken,
                // phone: phoneToken
            }
        }

        // save user in db
        await userModel.create(newUser)


        // email verification link
        const emailData = {
            from: USER,
            to: email,
            subject: "Verification Link",
            text: `${URL}/api/public/emailverify/${emailToken}`
        }

        sendEmail(emailData)


          // 8. Verification ke liye SMS data banate hain aur sendSMS function call karte hain.
        // const smsData = {
        //     body: `📲 Team Todo: Dear user, verify your phone by clicking the link: ${URL}/api/public/phoneverify/${phoneToken}. 
        //     If you didn't request this, ignore the message.`,
        //     to: phone
        // };
        
        // sendSMS(smsData);


        console.log(`${URL}/api/public/emailverify/${emailToken}`);
        // console.log(`${URL}/api/public/phoneverify/${phoneToken}`);
        
        res.status(200).json({msg: "You'll be registered as our new user, once u verify your email🙌"})

    } catch (error) {
        console.log(error);
        res.status(500).json({msg: error})
    }
})

router.get("/emailverify/:token", async (req: Request, res: Response): Promise<void> =>{

    try {

        // take token from url
        let token: string = req.params.token;
        
        // Compare URL token with userVerifyToken.email token
const user = await userModel.findOne({ "userVerifyToken.email": token });

if (!user) {
res.status(404).json({ msg: "Invalid token ❌" });
return
}

// If user already verified the link
if (user.userVerified.emailVerified) {
    res.status(200).json({ msg: "User email already verified" });
  return 
}

// userVerify token null and userVerified true
if (user) {  // ✅ Safe check before modifying user
  user.userVerified.emailVerified = true;
  user.userVerifyToken.emailToken = null;
  await user.save();  // ✅ Save the updated user
}

res.status(200).json({ msg: "User email verified successfully! ✅" });



    } catch (error) {
        console.log(error);
        res.status(500).json({msg: error})
    }

})

// router.get("/api/phoneverify/:token", async (req: Request, res: Response): Promise<void>=>{
//     try {
//         // take token from url
//         let token = req.params.token;

//         // compare token with emailToken
//         const user = await userModel.findOne({"userVerifyToken.phone": token})
//         if(!user){
//             res.status(200).json({msg: "Invalid Token"})
//             return;
//         }

//         //  check if user hasn't clicked link more than once
//         if(user.userVerified.phone === true){
//             res.status(200).json({msg: "User Phone Number Already Verified🙌"})
//             return
//         }

//         // change in db
//         user.userVerified.phone = true;
//         user.userVerifyToken.phone = null

//         // response 
//         res.status(200).json({msg: "user phone verified successfully!🙌"})
        
//     } catch (error) {
//         console.log(error);
//         res.status(500).json({msg: error})
//     }
// })


// login route
router.post("/usersignin", async (req: Request, res: Response): Promise<void>=>{
    try {

        // take input from user
        let {email, password} = req.body

        // check if email exists in db
        let checkUser = await userModel.findOne({email})
        if(!checkUser){
            res.status(200).json({msg: "email doesn't exists"})
            return
        }

        // compare password
        let pass = await bcrypt.compare(password, checkUser.password)
        if(!pass){
            res.status(200).json({msg: "invalid password"})
            return
        }

        // generate jwt token 
        let token = jwt.sign({id: checkUser._id}, KEY, {expiresIn: "30d"})

        
        res.status(200).json({msg: "User Logged in Successfully!🙌", token})
    } catch (error) {
        console.log(error);
        res.status(500).json({msg: error})
    }
})

router.post("/forgotpassword", async (req: Request, res: Response): Promise<void>=>{
    try {

        // take email from user
        let {email} = req.body;

        // check if exists in db
        let checkUser = await userModel.findOne({email})
        if(!checkUser){
            res.status(400).json({msg: "email not found"})
            return
        }

       // generate new password
       let newPass = Math.random().toString(36).substring(2)
       console.log(newPass);

    //   send pass on email

      const emailData = {
        from : USER,
        subject: "New Password",
        to: email,
         html: `<p>Your new password is: <strong>${newPass}</strong></p>`
      }
     sendEmail(emailData)

    //   hashs the new pass 
    let hashPass: string = await bcrypt.hash(newPass, 10)
    checkUser.password = hashPass; 
         
    await checkUser.save()
        
    res.status(200).json({msg: "New Password sent to ur email successfully!"})
    } catch (error) {
        console.log(error);
        res.status(500).json({msg: error})
    }


})

export default router