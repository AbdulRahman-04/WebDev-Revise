import mongoose, {Schema, Model, Document} from "mongoose";

interface ITodo extends Document {
    date: string,
    todoNo: number,
    todoTtitle: string,
    todoDescription: string,
    fileUpload: string
}

const todoSchema  = new Schema<ITodo>({
    date: {
        type: String,
        require: true,
    },
    todoNo: {
        type: Number,
        require: true
    },
    todoTtitle: {
        type: String,
        require: true
    },
    todoDescription: {
        type: String,
        require: true
    },
    // fileUpload: {
    //     type: String
    // }
}, {
    timestamps: true
})

const todoModel : Model<ITodo> = mongoose.model<ITodo>("todos", todoSchema, "todos")

export default todoModel  