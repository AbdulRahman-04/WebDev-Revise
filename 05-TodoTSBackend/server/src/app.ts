import express , {Request, Response, Application} from "express"
import config from "config"

const app : Application = express()

const PORT : string = config.get<string>("PORT")

app.use(express.json())

// get route 
app.get("/", async (req: Request, res: Response)=>{
    try {

        res.status(200).json({msg: "HELLO WORLD"})
        
    } catch (error) {
        console.log(error);
        res.status(500).json({msg: error})
    }
})


// server start 
app.listen(Number(PORT), ()=>{
    console.log(`Your server is live at port ${PORT}` );
    
})