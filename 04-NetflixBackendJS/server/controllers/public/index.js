import express from "express";
import config from "config";
import bcrypt from "bcrypt";
import jwt from "jsonwebtoken";
import adminModel from "../../models/ADMIN/Admin.js";
import sendEmail from "../../utils/sendEmail.js";
import sendSMS from "../../utils/sendSMS.js";

const router = express.Router()

const URL = config.get("URL");
const JWT_KEY = config.get("JWT_KEY");

router.post("/adminsignup", async (req, res) => {
  try {
    //  take input from admin
    let { username, email, password, age, phone } = req.body;

    // duplicate check
    let admin = await adminModel.findOne({ email });
    if (admin) {
      return res.status(200).json({ msg: "Admin already exists! Please login" });
    }

    // hash the password
    let hashPass = await bcrypt.hash(password, 10);

    // generate two random tokens for email and phone
    let emailToken = Math.random().toString(36).substring(2);
    let phoneToken = Math.random().toString(36).substring(2);

    // store the admin info in new obj and push into db
    let newAdmin = {
      username,
      email,
      password: hashPass,
      phone,
      age,
      adminVerifyToken: {
        email: emailToken,
        phone: phoneToken,
      },
    };

    await adminModel.create(newAdmin);

    // send verification link to email and phone number
    let emailData = {
        from: "Team Netflix",
        to: email,
        subject: "Email Verification",
        html: `
        <div style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 10px; background: #f9f9f9;">
            <div style="text-align: center; padding: 20px 0;">
                <img src="https://upload.wikimedia.org/wikipedia/commons/7/7a/Logonetflix.png" alt="Netflix Logo" width="150" />
                <h2 style="color: #e50914;">Welcome to Netflix</h2>
                <p style="font-size: 16px;">Dear Admin,</p>
                <p style="font-size: 16px;">Please verify your email by clicking the button below:</p>
            </div>
            <div style="text-align: center; margin: 20px 0;">
                <a href="${URL}/api/public/emailverify/${emailToken}" 
                   style="display: inline-block; background: #e50914; color: #fff; padding: 12px 30px; 
                   font-size: 16px; text-decoration: none; border-radius: 5px;">
                    Verify Email
                </a>
            </div>
            <div style="text-align: center; font-size: 14px; color: #555;">
                <p>If you did not sign up, you can safely ignore this email.</p>
                <p style="margin-top: 10px;">Thank you, <br> <strong>Team Netflix</strong></p>
            </div>
        </div>
        `,
        text: `${URL}/api/public/emailverify/${emailToken}`,
    };
    
    sendEmail(emailData);

    // send verification link to mobile phone number
    // let smsData = {
    //   body: `Dear Admin, please verify your phone here: ${URL}/api/public/phoneverify/${phoneToken}`,
    //   to: phone,
    // };

    // sendSMS(smsData);

    console.log(`${URL}/api/public/emailverify/${emailToken}`);
    // console.log(`${URL}/api/public/phoneverify/${phoneToken}`);

    res.status(200).json({
      msg: "You'll be registered as a new Netflix admin once you verify your email and mobile via the link provided on your email and phone number! 🙌",
    });

  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: error });
  }
});

router.get("/emailverify/:token", async (req, res) => {
  try {
    let token = req.params.token;

    let admin = await adminModel.findOne({ "adminVerifyToken.email": token });
    if (!admin) {
      return res.status(200).json({ msg: "Invalid token ❌" });
    }

    if (admin.adminVerified.email === true) {
      return res.status(200).json({ msg: "Admin email already verified 🙌" });
    }

    admin.adminVerified.email = true;
    admin.adminVerifyToken.email = null;

    await admin.save();

    res.status(200).json({ msg: `Email verified successfully! ✅` });

  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: error });
  }
});

router.get("/phoneverify/:token", async (req, res) => {
  try {
    let token = req.params.token;

    let admin = await adminModel.findOne({ "adminVerifyToken.phone": token });
    if (!admin) {
      return res.status(401).json({ msg: `Invalid token ❌` });
    }

    if (admin.adminVerified.phone === true) {
      return res.status(200).json({ msg: "Phone already verified 🙌" });
    }

    admin.adminVerified.phone = true;
    admin.adminVerifyToken.phone = null;

    await admin.save();

    res.status(200).json({ msg: `Phone verified ✅` });

  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: error });
  }
});

router.post("/adminsignin", async (req, res) => {
  try {
    let { email, password } = req.body;

    let admin = await adminModel.findOne({ email });
    if (!admin) {
      return res.status(200).json({ msg: "Email not found ❌" });
    }

    let checkPass = await bcrypt.compare(password, admin.password);
    if (!checkPass) {
      return res.status(404).json({ msg: `Invalid password` });
    }

    let token = jwt.sign({ _id: admin.id }, JWT_KEY, { expiresIn: "2d" });

    res.status(200).json({ msg: "Admin logged in successfully! ✅", token });

  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: error });
  }
});

export default router;
