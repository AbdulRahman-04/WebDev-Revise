import express from "express";
import moviesModel from "../../models/MOVIES/Movies.js";

const router = express.Router();


// 🛠️ Add new movie API
router.post("/add", async (req, res) => {
    try {
      const { movieName, releaseYear, movieGenre, rating, language, duration } = req.body;
  
      // ✅ Simple missing fields validation
      if (!movieName || !releaseYear || !movieGenre || !rating || !language || !duration) {
        return res.status(400).json({ msg: "Please fill all fields!" });
      }
  
      // ✅ Add movie to DB
      const newMovie = await moviesModel.create({
        movieName,
        releaseYear,
        movieGenre,
        rating,
        language,
        duration
      });
  
      res.status(201).json({ msg: "Movie added successfully!", newMovie });
  
    } catch (error) {
      console.log(error);
      res.status(500).json({ msg: "Internal server error" });
    }
  });

// ➤ GET ALL MOVIES
router.get("/getallmovies", async (req, res) => {
    try {
        let getAll = await moviesModel.find({});
        res.status(200).json({ msg: getAll });

    } catch (error) {
        console.log(error);
        res.status(500).json({ msg: error });
    }
});


// ➤ GET ONE MOVIE BY ID
router.get("/getonemovie/:id", async (req, res) => {
    try {
        let getOne = await moviesModel.findById(req.params.id);
        if (!getOne) {
            return res.status(404).json({ msg: "Movie not found!" });
        }
        res.status(200).json({ msg: getOne });

    } catch (error) {
        console.log(error);
        res.status(500).json({ msg: error });
    }
});


// ➤ UPDATE MOVIE BY ID
router.put("/updatemovie/:id", async (req, res) => {
    try {
        let userInput = req.body;
        let updatedMovie = await moviesModel.findByIdAndUpdate(
            req.params.id,
            { $set: userInput },
            { new: true }
        );

        if (!updatedMovie) {
            return res.status(404).json({ msg: "Movie not found!" });
        }

        res.status(200).json({ msg: "Movie updated successfully!🙌" });

    } catch (error) {
        console.log(error);
        res.status(500).json({ msg: error });
    }
});


// ➤ DELETE ONE MOVIE BY ID
router.delete("/deleteone/:id", async (req, res) => {
    try {
        let deletedMovie = await moviesModel.findByIdAndDelete(req.params.id);

        if (!deletedMovie) {
            return res.status(404).json({ msg: "Movie not found!" });
        }

        res.status(200).json({ msg: "Movie deleted successfully!🙌" });

    } catch (error) {
        console.log(error);
        res.status(500).json({ msg: error });
    }
});


// ➤ DELETE ALL MOVIES
router.delete("/deleteall", async (req, res) => {
    try {
        await moviesModel.deleteMany({});
        res.status(200).json({ msg: "All movies deleted successfully!🙌" });

    } catch (error) {
        console.log(error);
        res.status(500).json({ msg: error });
    }
});

export default router;
