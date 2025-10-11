import mongoose from "mongoose";

let seriesSchema = new mongoose.Schema({
    seriesName: {
        type: String,
        require: true
    },
    seriesEpisodes: {
        type: String,
         require: true
    },
    seriesSeasons: {
        type: String,
        require: true
    },
    releaseYear: {
        type: String,
        require: true
    },
    seriesGenre: {
        type: [String],
        require: true,
        enum: ["Horror, Comedy, Drama, Dramedy"]
    },
    rating :{
        type: String,
        require: true
    },
    language: {
        type: [String],
        enum: ["Hindi, English, Arabic, Kannada, Telugu, Marathi"],
        require: true
    }
}, {
    timestamps: true
})

let seriesModel =  mongoose.model("webseries", seriesSchema, "webseries")

export default seriesModel