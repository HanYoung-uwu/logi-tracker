import { NavigateFunction } from "react-router-dom";
import { API_URL_ROOT } from "../config/config";

interface ItemRecord {
    item: string
    location: string
    size: number
};

const fetchAllItems: (navigate: NavigateFunction) => Promise<Array<ItemRecord> | void> = async (navigate: NavigateFunction) => {
    let url = API_URL_ROOT + "/user/all_items"
    let res = await fetch(url);
    if (Math.floor(res.status / 100) != 2) {
        navigate("/login", {replace: true});
    } else {
        return await res.json();
    }
}

export {fetchAllItems}