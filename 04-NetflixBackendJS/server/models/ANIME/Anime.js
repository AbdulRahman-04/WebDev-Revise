import mongoose from "mongoose";

let animeSchema = new mongoose.Schema({
    animeName: {
        type: String,
        required: true
    },
    releaseYear: {
        type: String,
        required: true
    },
    animeGenre: {
        type: [String],
        required: true,
        enum: ["Action", "Adventure", "Fantasy", "Romance", "Slice of Life", "Sci-Fi", "Horror", "Comedy", "Drama", "Mystery", "Thriller"]
    },
    rating: {
        type: String,
        required: true
    },
    language: {
        type: [String],
        enum: ["Japanese", "English", "Hindi", "Spanish", "French", "German"],
        required: true
    }
}, {
    timestamps: true
});

let animeModel = mongoose.model("anime", animeSchema, "anime");

export default animeModel;
