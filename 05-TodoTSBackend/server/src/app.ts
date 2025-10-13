import express, {Request, Response, Application} from "express"
import config from "config"
import "./utils/dbConnect"

// security 
import helmet from "helmet"
import compression from "compression"
import rateLimit from "express-rate-limit"
import cors from "cors"

// public apis import 
import userRouter from "./controllers/public/users"

// middleware import 
import authMiddleware from "./middleware/auth"

// private apis import 
import todosRouter from "./controllers/private/todos"
import morgan from "morgan"

const app : Application = express()

const PORT: string = config.get<string>("PORT")

app.use(express.json())
app.use(express.urlencoded({ extended: true }));

// middlewares
app.use(helmet())
app.use(compression())
app.use(cors())

app.use(morgan("dev"))

// rate limit 
const limiter = rateLimit({
    windowMs: 15 * 60 *1000,
    max: 100,
    message: "Too many requests"
})
app.use(limiter)

// route 
app.get("/", async (req: Request, res: Response)=>{
    try {

        res.status(200).json({msg: "HELLO WORLD"})
        
    } catch (error) {
        console.log(error);
        res.status(500).json({msg: error})
    }
})

// fallback route 
app.use("*", (_req, res) => {
  res.status(404).json({ msg: "Route not found ❌" });
});


// public apis 
app.use("/api/public", userRouter)

// private apis 
app.use("/api/private", authMiddleware, todosRouter)


app.listen(Number(PORT), ()=>{
    console.log(`server is live at port ${PORT}`);
    
})