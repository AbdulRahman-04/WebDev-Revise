import express from "express";
import adminModel from "../../models/ADMIN/Admin.js";

const router = express.Router();

// ➤ GET ALL ADMINS
router.get("/getalladmins", async (req, res) => {
    try {
        let getAll = await adminModel.find({});
        res.status(200).json({ msg: getAll });

    } catch (error) {
        console.log(error);
        res.status(500).json({ msg: error });
    }
});

// ➤ GET ONE ADMIN BY ID
router.get("/getoneadmin/:id", async (req, res) => {
    try {
        let getOne = await adminModel.findById(req.params.id);

        if (!getOne) {
            return res.status(404).json({ msg: "Admin not found!" });
        }

        res.status(200).json({ msg: getOne });

    } catch (error) {
        console.log(error);
        res.status(500).json({ msg: error });
    }
});

// ➤ CREATE NEW ADMIN
router.post("/createadmin", async (req, res) => {
    try {
        let newAdmin = new adminModel(req.body);
        await newAdmin.save();

        res.status(201).json({ msg: "Admin created successfully! 🙌", admin: newAdmin });

    } catch (error) {
        console.log(error);
        res.status(500).json({ msg: error });
    }
});

// ➤ UPDATE ADMIN BY ID
router.put("/updateadmin/:id", async (req, res) => {
    try {
        let adminInput = req.body;
        let updatedAdmin = await adminModel.findByIdAndUpdate(
            req.params.id,
            { $set: adminInput },
            { new: true }
        );

        if (!updatedAdmin) {
            return res.status(404).json({ msg: "Admin not found!" });
        }

        res.status(200).json({ msg: "Admin updated successfully! 🙌", admin: updatedAdmin });

    } catch (error) {
        console.log(error);
        res.status(500).json({ msg: error });
    }
});

// ➤ DELETE ONE ADMIN BY ID
router.delete("/deleteone/:id", async (req, res) => {
    try {
        let deletedAdmin = await adminModel.findByIdAndDelete(req.params.id);

        if (!deletedAdmin) {
            return res.status(404).json({ msg: "Admin not found!" });
        }

        res.status(200).json({ msg: "Admin deleted successfully! 🙌" });

    } catch (error) {
        console.log(error);
        res.status(500).json({ msg: error });
    }
});

// ➤ DELETE ALL ADMINS
router.delete("/deleteall", async (req, res) => {
    try {
        await adminModel.deleteMany({});
        res.status(200).json({ msg: "All admins deleted successfully! 🙌" });

    } catch (error) {
        console.log(error);
        res.status(500).json({ msg: error });
    }
});

export default router;
