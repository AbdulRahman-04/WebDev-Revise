import express from "express";
import config from "config";
import "./utils/dbConnect.js";

// public routes
import adminPublicRouter from "./controllers/public/index.js" 

// auth middleware
import authMiddleware from "./middleware/auth.js";

// private routes
import movieRouter from "./controllers/Movies/index.js";
import seriesRouter from "./controllers/Webseries/index.js";
import animeRouter from "./controllers/Anime/index.js";
import adminPrivate from "./controllers/Admins/index.js"

// rate limit
import ratelimit from "express-rate-limit";

const app = express();
const PORT = config.get("PORT");

app.use(express.json());

// Rate limiter configuration
let limiter = ratelimit({
    windowMs: 10 * 60 * 1000,  // 10 min
    limit: 50,
    standardHeaders: true,
    legacyHeaders: false,
    message: "Cannot send request! Wait for the server to respond!",
    statusCode: 429
});

app.get("/", (req, res) => {
    try {
        res.status(200).json({ msg: 'Welcome to Netflix Admin Panel' });  // Updated message
    } catch (error) {
        console.log(error);
        res.status(401).json({ msg: error });
    }
});


app.use("/api/public", adminPublicRouter);



// Private APIs
app.use("/api/admins", authMiddleware, adminPrivate)
app.use("/api/movies", authMiddleware, movieRouter);
app.use("/api/series",authMiddleware, seriesRouter);
app.use("/api/anime",authMiddleware, animeRouter);


app.listen(PORT, () => {
    console.log(`Your Netflix Admin Panel is running live at port ${PORT}`);  // Updated message
});
