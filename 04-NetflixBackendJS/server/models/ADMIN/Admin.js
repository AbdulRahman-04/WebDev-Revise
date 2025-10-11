import mongoose from "mongoose";

let adminSchema = new mongoose.Schema({
    username: {
        type: String,
        required: true,
        maxlength: 50,
        minlength: 5
    },
    email: {
        type: String,
        required: true,
        unique: true
    },
    password: {
        type: String,
        required: true
    },
    phone: {
        type: String,
        required: true
    },
    age: {
        type: Number,
        required: true
    },
    isAdmin: {
        type: Boolean,
        default: true   // Admin flag
    },
    adminVerified: {
        email: {
            type: Boolean,
            default: false
        },
        phone: {
            type: Boolean,
            default: false
        }
    },
    adminVerifyToken: {
        email: {
            type: String
        },
        phone: {
            type: String
        }
    }
},
{
    timestamps: true
});

let adminModel = mongoose.model("admins", adminSchema, "admins");

export default adminModel;
