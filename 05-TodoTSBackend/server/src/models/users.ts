import mongoose, {Schema, Document, Model} from "mongoose";

// create interface 
interface IUser extends Document{
    userName: string,
    email: string,
    password: string,
    age: number,
    userVerified: {
        emailVerified: boolean | null,
        // phoneVerified: boolean | null,
    },
    userVerifyToken : {
        emailToken: string | null,
        // phoneToken : string | null
    }
}

const UserSchema = new Schema<IUser>({
      userName: {
        type: String,
        required: true,
        maxlength: 70,
        minlength: 10
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
    age: {
        type: Number,
        required: true
    },
    userVerified: {
        emailVerified: {
            type: Boolean,
            default: false
        },
        phoneVerified: {
            type: Boolean,
            default: false
        }
    },
    userVerifyToken: {
        email: {
            type: String,  // Changed to String
            default: null
        },
        phone: {
            type: String,  // Changed to String
            default: null
        }
    }
}, {
    timestamps: true
})


const userModel: Model<IUser> = mongoose.model<IUser>("users", UserSchema, "users")

export {userModel, IUser}