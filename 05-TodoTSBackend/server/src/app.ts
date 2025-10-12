import express, {Request, Response, Application} from "express"
import config from "config"
import "./utils/dbConnect"

const app : Application = express()

const PORT: string = config.get<string>("PORT")

app.use(express.json())

// route 
app.get("/", async (req: Request, res: Response)=>{
    try {

        res.status(200).json({msg: "HELLO WORLD"})
        
    } catch (error) {
        console.log(error);
        res.status(500).json({msg: error})
    }
})


app.listen(Number(PORT), ()=>{
    console.log(`server is live at port ${PORT}`);
    
})