import express, { Request, Response } from "express";
import todoModel from "../../models/todos";

const router = express.Router();

router.post("/addtodo", async (req: Request, res: Response): Promise<void> => {
  try {
    const { date, todoNo, todoTtitle, todoDescription, fileUpload } = req.body;

    if (!date || !todoNo || !todoTtitle || !todoDescription) {
      res.status(400).json({ msg: "Missing required fields ❌" });
      return;
    }

    const newTodo = new todoModel({
      date,
      todoNo,
      todoTtitle,
      todoDescription,
    });

    await newTodo.save();
    res.status(201).json({ msg: "Todo added ✅", todo: newTodo });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: "Server error ❌" });
  }
});

router.get("/alltodos", async (_req: Request, res: Response): Promise<void> => {
  try {
    const todos = await todoModel.find({});
    res.status(200).json({ todos });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: "Server error ❌" });
  }
});

router.get("/getone/:id", async (req: Request, res: Response): Promise<void> => {
  try {
    const todo = await todoModel.findById(req.params.id);
    if (!todo) {
      res.status(404).json({ msg: "Todo not found ❌" });
      return;
    }
    res.status(200).json({ todo });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: "Server error ❌" });
  }
});

router.put("/editone/:id", async (req: Request, res: Response): Promise<void> => {
  try {
    const { date, todoNo, todoTtitle, todoDescription, fileUpload } = req.body;

    const updatedTodo = await todoModel.findByIdAndUpdate(
      req.params.id,
      { date, todoNo, todoTtitle, todoDescription, fileUpload },
      { new: true }
    );

    if (!updatedTodo) {
      res.status(404).json({ msg: "Todo not found ❌" });
      return;
    }

    res.status(200).json({ msg: "Todo updated ✅", updatedTodo });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: "Server error ❌" });
  }
});

router.delete("/deleteone/:id", async (req: Request, res: Response): Promise<void> => {
  try {
    const deletedTodo = await todoModel.findByIdAndDelete(req.params.id);
    if (!deletedTodo) {
      res.status(404).json({ msg: "Todo not found ❌" });
      return;
    }
    res.status(200).json({ msg: "Todo deleted ✅" });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: "Server error ❌" });
  }
});

router.delete("/deleteall", async (_req: Request, res: Response): Promise<void> => {
  try {
    await todoModel.deleteMany({});
    res.status(200).json({ msg: "All todos deleted ✅" });
  } catch (error) {
    console.log(error);
    res.status(500).json({ msg: "Server error ❌" });
  }
});

export default router