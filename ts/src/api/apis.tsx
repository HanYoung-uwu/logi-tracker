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
    time: Date
    code: string
}

interface HistoryRecord {
    Location: string
    User: string
    Time: Date
    Action: string
    ItemType: string
    Size: number
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

const addItem: (item: string, size: number, location: string) => Promise<boolean> = async (item, size, location) => {
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
        return false;
    };
    return true;
};

const deleteItem: (item: string, location: string) => void = async (item, location) => {
    let url = API_URL_ROOT + "/user/delete_item";
    let headers = {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    };
    let res = await fetch(url, {
        method: "POST",
        headers: headers,
        body: JSON.stringify({
            item: item,
            location: location
        })
    });
    if (Math.floor(res.status / 100) != 2) {
        return null;
    };
};

const setItem: (item: string, size: number, location: string) => Promise<boolean> = async (item, size, location) => {
    let url = API_URL_ROOT + "/user/set_item";
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
        return false;
    };
    return true;
};

const fetchHistory: (navigate: NavigateFunction) => Promise<Array<HistoryRecord>> = async (navigate) => {
    let url = API_URL_ROOT + "/user/history";
    let res = await fetch(url);
    if (Math.floor(res.status / 100) != 2) {
        navigate("/login", { replace: true });
        return [];
    };
    let result = await res.json();
    let returnValue = [];
    for (let i = 0; i < result.length; i++) {
        let action = '';
        switch (result[i].Action) {
            case 0:
                action = "add"; break;
            case 1:
                action = "retrieve"; break;
            case 2:
                action = "delete"; break;
            case 3:
                action = "set"; break;
        }
        result[i].Action = action;
        result[i].Time = new Date(result[i].Time);
    }
    return result;
}

export { fetchAllItems, fetchAccountInfo, fetchClanInviteLink, addItem, fetchLocations, deleteItem, setItem, fetchHistory };
export type { Location, HistoryRecord };
