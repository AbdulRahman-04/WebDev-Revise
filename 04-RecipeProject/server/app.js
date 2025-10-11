import express from "express"
import config from "config"
// db connect
import "./utils/dbConnect.js"

const app = express()

const PORT = config.get("PORT")

app.use(express.json())

app.get("/", async (req, res)=>{
    try {

        res.status(200).json({msg: "Hello World!"})
        
    } catch (error) {
        console.log(error);
        res.status(500).json({msg: error})
        
    }
})

app.listen(PORT, ()=>{
    console.log(`Your server is running live at port ${PORT} `);
    
})