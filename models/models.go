package models

import (
	"database/sql"
	"fmt"
)

type Daemon struct {
	ID         int            `db:"id"`
	WalletUser string         `db:"wallet_usr"`
	WalletPass string         `db:"wallet_pass"`
	WalletPort int            `db:"wallet_port"`
	Folder     string         `db:"folder"`
	NodeID     int            `db:"node_id"`
	CoinID     int            `db:"coin_id"`
	Conf       string         `db:"conf"`
	IP         string         `db:"ip"`
	MnPort     int            `db:"mn_port"`
	PassPhrase sql.NullString `db:"wallet_passphrase"`
}

func (d *Daemon) ToString() string {
	return fmt.Sprintf("[%s, %s, %d, %s, %d, %d, %s, %s, %d]", d.WalletUser, d.WalletPass, d.WalletPort, d.Folder, d.NodeID, d.CoinID, d.Conf, d.IP, d.MnPort)
}
