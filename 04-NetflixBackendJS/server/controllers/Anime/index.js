import express from "express";
import animeModel from "../../models/ANIME/Anime.js";

const router = express.Router();

// 🛠️ Add new anime API
router.post("/add", async (req, res) => {
  try {
    const {
      animeName,
      releaseYear,
      genre,
      rating,
      language,
      episodes,
      duration,
    } = req.body;

    // ✅ Simple missing fields validation
    if (
      !animeName ||
      !releaseYear ||
      !genre ||
      !rating ||
      !language ||
      !episodes ||
      !duration
    ) {
      return res.status(400).json({ msg: "Please fill all fields!" });
    }

    // ✅ Add anime to DB
    const newAnime = await animeModel.create({
      animeName,
      releaseYear,
      genre,
      rating,
      language,
      episodes,
      duration,
    });

    res.status(201).json({ msg: "Anime added successfully!", newAnime });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: "Internal server error" });
  }
});

// ➤ GET ALL ANIME
// ✔️ Fetches all anime documents from the database
router.get("/getallanime", async (req, res) => {
  try {
    let getAll = await animeModel.find({});
    res.status(200).json({ msg: getAll });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: error });
  }
});

// ➤ GET ONE ANIME BY ID
// ✔️ Fetches a single anime by ID
router.get("/getoneanime/:id", async (req, res) => {
  try {
    let paramsId = req.params.id;

    // ❌ Issue: Incorrect query filter (fix below)
    // let getOne = await animeModel.findOne({ paramsId }); ❌ WRONG

    // ✅ Correct filter query
    let getOne = await animeModel.findById(paramsId);

    if (!getOne) {
      return res.status(404).json({ msg: "Anime not found!" });
    }

    res.status(200).json({ msg: getOne });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: error });
  }
});

// ➤ UPDATE ANIME BY ID
// ✔️ Updates the anime document by ID
router.put("/updateanime/:id", async (req, res) => {
  try {
    let userInput = req.body;

    // ✅ Update anime by ID with new data
    let updatedAnime = await animeModel.findByIdAndUpdate(
      req.params.id,
      { $set: userInput },
      { new: true }
    );

    if (!updatedAnime) {
      return res.status(404).json({ msg: "Anime not found!" });
    }

    res
      .status(200)
      .json({ msg: "Anime edited successfully! 🙌", anime: updatedAnime });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: error });
  }
});

// ➤ DELETE ONE ANIME BY ID
// ✔️ Deletes a single anime by ID
router.delete("/deleteone/:id", async (req, res) => {
  try {
    let deletedAnime = await animeModel.findByIdAndDelete(req.params.id);

    if (!deletedAnime) {
      return res.status(404).json({ msg: "Anime not found!" });
    }

    res.status(200).json({ msg: "Anime deleted successfully! 🙌" });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: error });
  }
});

// ➤ DELETE ALL ANIME
// ✔️ Deletes all anime documents from the database
router.delete("/deleteall", async (req, res) => {
  try {
    await animeModel.deleteMany({});
    res.status(200).json({ msg: "All anime deleted successfully! 🙌" });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: error });
  }
});

export default router;
