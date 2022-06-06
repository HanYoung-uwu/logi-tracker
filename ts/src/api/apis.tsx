import { NavigateFunction } from "react-router-dom";
import { API_URL_ROOT } from "../config/config";

interface ItemRecord {
    item: string
    location: string
    size: number
};

interface AccountInfo {
    Name: string
    Clan: string
    Permission: 0 | 1 | 2 // 0 admin, 1 clan admin, 2 clan man
}

const fetchAllItems: (navigate: NavigateFunction) => Promise<Array<ItemRecord> | void> = async (navigate: NavigateFunction) => {
    let url = API_URL_ROOT + "/user/all_items"
    let res = await fetch(url);
    if (Math.floor(res.status / 100) != 2) {
        navigate("/login", {replace: true});
    } else {
        return await res.json();
    }
}

const fetchAccountInfo: () => Promise<AccountInfo | null> = async () => {
    let url = API_URL_ROOT + "/user/info"
    let res = await fetch(url);
    if (Math.floor(res.status / 100) != 2) {
        return null;
    } else {
        return await res.json();
    }
}

export {fetchAllItems, fetchAccountInfo}