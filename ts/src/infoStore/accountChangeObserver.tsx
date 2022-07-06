import { useContext, useEffect } from "react";
import { observer } from "mobx-react-lite";
import { AccountInfoContext, AccountInfo } from "./accountInfoStore";
import { LocationInfoContext, StockpileInfo } from "./stockPileInfoStore";
import { ItemsStoreContext } from "./itemsStore";

const AccountChangeObserver = observer(() => {
    const accountInfo = useContext(AccountInfoContext);
    const stockpileInfo = useContext(LocationInfoContext);
    const itemsStore = useContext(ItemsStoreContext);

    useEffect(() => {
        if (accountInfo.name !== "") {
            stockpileInfo.refetchInfo();
            itemsStore.refetchInfo();
        }
    }, [accountInfo.name]);
    return null;
});

export default AccountChangeObserver;