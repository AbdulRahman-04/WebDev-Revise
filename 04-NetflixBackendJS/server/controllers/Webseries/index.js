import express from "express";
import seriesModel from "../../models/WEBSERIES/Webseries.js";

const router = express.Router();

// 🛠️ Add new webseries API
router.post("/add", async (req, res) => {
  try {
    const {
      seriesName,
      releaseYear,
      genre,
      rating,
      language,
      seasons,
      episodes,
    } = req.body;

    // ✅ Simple missing fields validation
    if (
      !seriesName ||
      !releaseYear ||
      !genre ||
      !rating ||
      !language ||
      !seasons ||
      !episodes
    ) {
      return res.status(400).json({ msg: "Please fill all fields!" });
    }

    // ✅ Add webseries to DB
    const newSeries = await webseriesModel.create({
      seriesName,
      releaseYear,
      genre,
      rating,
      language,
      seasons,
      episodes,
    });

    res.status(201).json({ msg: "Webseries added successfully!", newSeries });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: "Internal server error" });
  }
});

// ➤ GET ALL SERIES
router.get("/getallseries", async (req, res) => {
  try {
    let getAll = await seriesModel.find({});
    res.status(200).json({ msg: getAll });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: error });
  }
});

// ➤ GET ONE SERIES BY ID
router.get("/getoneseries/:id", async (req, res) => {
  try {
    let getOne = await seriesModel.findById(req.params.id);
    if (!getOne) {
      return res.status(404).json({ msg: "Series not found!" });
    }
    res.status(200).json({ msg: getOne });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: error });
  }
});

// ➤ UPDATE SERIES BY ID
router.put("/updateseries/:id", async (req, res) => {
  try {
    let userInput = req.body;
    let updatedSeries = await seriesModel.findByIdAndUpdate(
      req.params.id,
      { $set: userInput },
      { new: true }
    );

    if (!updatedSeries) {
      return res.status(404).json({ msg: "Series not found!" });
    }

    res.status(200).json({ msg: "Series updated successfully!🙌" });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: error });
  }
});

// ➤ DELETE ONE SERIES BY ID
router.delete("/deleteone/:id", async (req, res) => {
  try {
    let deletedSeries = await seriesModel.findByIdAndDelete(req.params.id);

    if (!deletedSeries) {
      return res.status(404).json({ msg: "Series not found!" });
    }

    res.status(200).json({ msg: "Series deleted successfully!🙌" });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: error });
  }
});

// ➤ DELETE ALL SERIES
router.delete("/deleteall", async (req, res) => {
  try {
    await seriesModel.deleteMany({});
    res.status(200).json({ msg: "All series deleted successfully!🙌" });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: error });
  }
});

export default router;
