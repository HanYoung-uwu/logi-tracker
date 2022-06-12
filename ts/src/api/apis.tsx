import { Item } from "framer-motion/types/components/Reorder/Item";
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

interface Location {
    location: string
	time:     Date
	code:     string
}

const fetchAllItems: (navigate: NavigateFunction) => Promise<Array<ItemRecord> | void> = async (navigate: NavigateFunction) => {
    let url = API_URL_ROOT + "/user/all_items"
    let res = await fetch(url);
    if (Math.floor(res.status / 100) != 2) {
        navigate("/login", { replace: true });
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

const fetchClanInviteLink: () => Promise<string | null> = async () => {
    let url = API_URL_ROOT + "/clan/invitation";
    let res = await fetch(url);
    if (Math.floor(res.status / 100) != 2) {
        return null;
    } else {
        return (await res.json())["token"];
    }
};

const fetchLocations: () => Promise<Array<Location> | null> = async () => {
    let url = API_URL_ROOT + "/user/all_stockpiles";
    let res = await fetch(url);
    if (Math.floor(res.status / 100) != 2) {
        return null;
    } else {
        let json = await res.json();
        return json.map((location: any) => {
            return {
                location: location.Location,
                code: location.Code,
                time: new Date(location.Time)
            };
        });
    }
}

const addItem: (item: string, size: number, location: string) => void = async (item, size, location) => {
    let url = API_URL_ROOT + "/user/update_item";
    let headers = {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    };
    let res = await fetch(url, {
        method: "POST",
        headers: headers,
        body: JSON.stringify({
            item: item,
            location: location,
            size: size
        })
    });
    if (Math.floor(res.status / 100) != 2) {
        return null;
    };
};

export { fetchAllItems, fetchAccountInfo, fetchClanInviteLink, addItem, fetchLocations };
export type { Location };
