import mongoose from "mongoose";

let moviesSchema = new mongoose.Schema({
    movieName: {
        type: String,
        required: true
    },
    releaseYear: {
        type: String,
        required: true
    },
    movieGenre: {
        type: [String],
        required: true,
        enum: ["Horror", "Comedy", "Drama", "Thriller", "Action", "Romance", "Sci-Fi", "Documentary"]
    },
    rating: {
        type: String,
        required: true
    },
    language: {
        type: [String],
        enum: ["Hindi", "English", "Arabic", "Kannada", "Telugu", "Marathi", "Tamil", "Spanish"],
        required: true
    },
    duration: {
        type: String,  // e.g., "2h 15m"
        required: true
    }
}, {
    timestamps: true
});

let moviesModel = mongoose.model("movies", moviesSchema, "movies");

export default moviesModel;
