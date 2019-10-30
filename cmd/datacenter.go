package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hichikaw/dosanco/model"
)

func showDataCenter(cmd *cobra.Command, args []string) {
	url := Conf.APIServer.URL + "/datacenter"
	if len(args) > 0 {
		// show specified datacenter
		body, err := sendRequest("GET", url+"/name/"+args[0], []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		dc := new(model.DataCenter)
		if err = json.Unmarshal(body, dc); err != nil {
			fmt.Println("response parse error", err)
		}
		url = url + "/" + strconv.Itoa(int(dc.ID))
		body, err = sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		data := new(model.DataCenter)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("response parse error", err)
			return
		}
		data.Write(cmd.Flag("output").Value.String())
	} else {
		// show all datacenters
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		data := new(model.DataCenters)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("json unmarshall error:", err)
			return
		}
		data.Write(cmd.Flag("output").Value.String())
	}
}

func showDataCenterFloor(cmd *cobra.Command, args []string) {
	url := Conf.APIServer.URL + "/datacenter/floor"
	dcName := cmd.Flag("dc").Value.String()
	if len(args) > 0 {
		// show specified datacenter floor
		if dcName == "" {
			fmt.Println("datacenter name is required for showing specific floor")
			return
		}
		body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/name/"+dcName, []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		dc := new(model.DataCenter)
		if err = json.Unmarshal(body, dc); err != nil {
			fmt.Println("response parse error", err)
		}
		dcID := strconv.Itoa(int(dc.ID))
		body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		data := new(model.Floors)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("response parse error:", err)
			return
		}
		floor := model.Floor{}
		for _, flr := range *data {
			if flr.Name == args[0] {
				floor = flr
			}
		}
		if floor.ID != 0 {
			floor.Write(cmd.Flag("output").Value.String())
		} else {
			fmt.Println("floor not found")
		}
	} else {
		if dcName != "" {
			// show all datacenter floors
			url = Conf.APIServer.URL + "/datacenter/name/" + dcName
			body, err := sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err)
				return
			}
			dc := new(model.DataCenter)
			if err = json.Unmarshal(body, dc); err != nil {
				fmt.Println("response parse error", err)
			}
			dcID := strconv.Itoa(int(dc.ID))
			body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
			if err != nil {
				fmt.Println(err)
				return
			}
			data := new(model.Floors)
			if err := json.Unmarshal(body, data); err != nil {
				fmt.Println("response parse error:", err)
				return
			}
			outputModel := model.Floors{}
			for _, f := range *data {
				f.DataCenter = *dc
				outputModel = append(outputModel, f)
			}
			outputModel.Write(cmd.Flag("output").Value.String())
		} else {
			// show all datacenter floors
			body, err := sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err)
				return
			}
			data := new(model.Floors)
			if err := json.Unmarshal(body, data); err != nil {
				fmt.Println("response parse error:", err)
				return
			}
			// get datacenters
			dcs := new(model.DataCenters)
			url = Conf.APIServer.URL + "/datacenter"
			body, err = sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if err := json.Unmarshal(body, dcs); err != nil {
				fmt.Println("response parse error:", err)
				return
			}
			outputModel := model.Floors{}
			for _, f := range *data {
				dc, err := dcs.Take(f.DataCenterID)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				f.DataCenter = *dc
				outputModel = append(outputModel, f)
			}

			outputModel.Write(cmd.Flag("output").Value.String())
		}
	}
}

func showDataCenterHall(cmd *cobra.Command, args []string) {
	url := Conf.APIServer.URL + "/datacenter/hall"
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	if len(args) > 0 {
		// show specified datacenter floor
		if dcName == "" {
			fmt.Println("datacenter name is required for showing specific hall")
			return
		}
		if floorName == "" {
			fmt.Println("floor name is required for showing specific hall")
			return
		}
		body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/name/"+dcName, []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		dc := new(model.DataCenter)
		if err = json.Unmarshal(body, dc); err != nil {
			fmt.Println("response parse error", err)
		}
		dcID := strconv.Itoa(int(dc.ID))
		body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		data := new(model.Floors)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("response parse error:", err)
			return
		}
		floor := model.Floor{}
		for _, flr := range *data {
			if flr.Name == floorName {
				floor = flr
			}
		}
		if floor.ID == 0 {
			fmt.Println("floor not found")
		}

		fID := strconv.Itoa(int(floor.ID))
		f := new(model.Floor)
		body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/floor/"+fID, []byte{})
		if err != nil {
			fmt.Println(err)
			return
		}
		if err = json.Unmarshal(body, f); err != nil {
			fmt.Println(err)
			return
		}
		for _, h := range f.Halls {
			if h.Name == args[0] {
				h.Write(cmd.Flag("output").Value.String())
			}
		}
	} else {
		if dcName != "" {
			// show all halls of specified datacenter
			url = Conf.APIServer.URL + "/datacenter/name/" + dcName
			body, err := sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err)
				return
			}
			dc := new(model.DataCenter)
			if err = json.Unmarshal(body, dc); err != nil {
				fmt.Println("response parse error", err)
			}
			dcID := strconv.Itoa(int(dc.ID))
			body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
			if err != nil {
				fmt.Println(err)
				return
			}
			data := new(model.Floors)
			if err := json.Unmarshal(body, data); err != nil {
				fmt.Println("response parse error:", err)
				return
			}
			if floorName == "" {
				halls := model.Halls{}
				for _, f := range *data {
					floor := new(model.Floor)
					fID := strconv.Itoa(int(f.ID))
					body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/floor/"+fID, []byte{})
					if err != nil {
						fmt.Println(err)
						return
					}
					if err = json.Unmarshal(body, floor); err != nil {
						fmt.Println(err)
						return
					}
					if len(floor.Halls) > 0 {
						for _, h := range floor.Halls {
							h.Floor = *floor
							halls = append(halls, h)
						}
					}
				}
				halls.Write(cmd.Flag("output").Value.String())
			} else {
				floor := new(model.Floor)
				for _, f := range *data {
					if f.Name == floorName {
						floor = &f
						break
					}
				}
				fID := strconv.Itoa(int(floor.ID))
				body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/floor/"+fID, []byte{})
				if err != nil {
					fmt.Println(err)
					return
				}
				if err = json.Unmarshal(body, floor); err != nil {
					fmt.Println(err)
					return
				}
				outputModel := model.Halls{}
				for _, h := range floor.Halls {
					h.Floor = *floor
					outputModel = append(outputModel, h)
				}
				outputModel.Write(cmd.Flag("output").Value.String())
			}
		} else {
			// show all datacenter halls
			body, err := sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err)
				return
			}
			data := new(model.Halls)
			if err := json.Unmarshal(body, data); err != nil {
				fmt.Println("response parse error:", err)
				return
			}
			// get all floors
			url = Conf.APIServer.URL + "/datacenter/floor"
			body, err = sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err)
				return
			}
			floors := new(model.Floors)
			if err := json.Unmarshal(body, floors); err != nil {
				fmt.Println("response parse error:", err)
				return
			}
			outputModel := model.Halls{}
			for _, h := range *data {
				floor, err := floors.Take(h.FloorID)
				if err != nil {
					fmt.Println(err)
					return
				}
				h.Floor = *floor
				outputModel = append(outputModel, h)
			}
			outputModel.Write(cmd.Flag("output").Value.String())
		}
	}
}

func showRackRow(cmd *cobra.Command, args []string) {
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	if len(args) > 0 {
		// show specified row
		url := Conf.APIServer.URL + "/datacenter/row?dc=" + dcName + "&floor=" + floorName + "&hall=" + hallName + "&name=" + args[0]
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := new(model.RackRows)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("parse response error")
			return
		}
		if len(*data) > 1 {
			fmt.Println("multiple row found")
			return
		}
		for _, r := range *data {
			// get hall
			hallID := strconv.Itoa(int(r.HallID))
			url = Conf.APIServer.URL + "/datacenter/hall/" + hallID
			body, err := sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if err := json.Unmarshal(body, &r.Hall); err != nil {
				fmt.Println("parse response error")
				return
			}
			r.Write(cmd.Flag("output").Value.String())
			break
		}
	} else {
		// show list of rows
		url := Conf.APIServer.URL + "/datacenter/row?dc=" + dcName
		if floorName != "" {
			url = url + "&floor=" + floorName
		}
		if hallName != "" {
			url = url + "&hall=" + hallName
		}
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := new(model.RackRows)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("parse response error")
			return
		}

		outputModel := model.RackRows{}
		// get halls
		url = Conf.APIServer.URL + "/datacenter/hall?dc=" + dcName
		if floorName != "" {
			url = url + "&floor=" + floorName
		}
		body, err = sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		halls := new(model.Halls)
		if err := json.Unmarshal(body, halls); err != nil {
			fmt.Println("parse response error")
			return
		}

		// get floor
		fm := make(map[uint]struct{})
		floorSlice := []model.Floor{}
		for _, hall := range *halls {
			fm[hall.FloorID] = struct{}{}
		}
		for id := range fm {
			floor, err := getFloor(id)
			if err != nil {
				fmt.Println("get floor error")
			}
			floorSlice = append(floorSlice, *floor)
		}
		floors := model.Floors(floorSlice)

		// get datacenter
		dcm := make(map[uint]struct{})
		dcSlice := []model.DataCenter{}
		for _, floor := range floors {
			dcm[floor.DataCenterID] = struct{}{}
		}
		for id := range dcm {
			dc, err := getDataCenter(id)
			if err != nil {
				fmt.Println("get floor error")
			}
			dcSlice = append(dcSlice, *dc)
		}
		dcs := model.DataCenters(dcSlice)

		for _, r := range *data {
			hall, err := halls.Take(r.HallID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			f, _ := floors.Take(hall.FloorID)
			d, _ := dcs.Take(f.DataCenterID)
			f.DataCenter = *d
			hall.Floor = *f
			r.Hall = *hall
			outputModel = append(outputModel, r)
		}
		outputModel.Write(cmd.Flag("output").Value.String())
	}
}

func showRack(cmd *cobra.Command, args []string) {
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	rowName := cmd.Flag("row").Value.String()
	rowPduName := cmd.Flag("row-pdu").Value.String()
	if len(args) > 0 {
		// show specified rack
		if dcName == "" {
			fmt.Println("datacenter name is required")
			return
		}
		if floorName == "" {
			fmt.Println("floor name is required")
			return
		}
		if hallName == "" {
			fmt.Println("hall name is required")
			return
		}
		if rowName == "" {
			fmt.Println("row name is required")
			return
		}
		if rowPduName == "" {
			fmt.Println("row-pdu name is ignored")
		}
		url := Conf.APIServer.URL + "/datacenter/rack?dc=" + dcName + "&floor=" + floorName + "&hall=" + hallName + "&row=" + rowName + "&name=" + args[0]
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := new(model.Racks)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("parse response error")
			return
		}
		for _, r := range *data {
			rowID := strconv.Itoa(int(r.RowID))
			url = Conf.APIServer.URL + "/datacenter/row/" + rowID
			body, err := sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if err := json.Unmarshal(body, &r.RackRow); err != nil {
				fmt.Println("parse response error")
				return
			}
			r.Write(cmd.Flag("output").Value.String())
			break
		}
	} else {
		// show list of racks
		url := Conf.APIServer.URL + "/datacenter/rack?dc=" + dcName
		if floorName != "" {
			url = url + "&floor=" + floorName
		}
		if hallName != "" {
			url = url + "&hall=" + hallName
		}
		if rowName != "" {
			url = url + "&row=" + rowName
		}
		if rowPduName != "" {
			url = url + "&pdu=" + rowPduName
		}
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := new(model.Racks)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("parse response error")
			return
		}
		// get rack row
		body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/row?dc="+dcName, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		rows := new(model.RackRows)
		if err := json.Unmarshal(body, rows); err != nil {
			fmt.Println("parse response error")
			return
		}

		// get hall
		m := make(map[uint]struct{})
		hallSlice := []model.Hall{}
		for _, row := range *rows {
			m[row.HallID] = struct{}{}
		}
		for id := range m {
			hall, err := getHall(id)
			if err != nil {
				fmt.Println("get hall error")
			}
			hallSlice = append(hallSlice, *hall)
		}
		halls := model.Halls(hallSlice)

		// get floor
		fm := make(map[uint]struct{})
		floorSlice := []model.Floor{}
		for _, hall := range halls {
			fm[hall.FloorID] = struct{}{}
		}
		for id := range fm {
			floor, err := getFloor(id)
			if err != nil {
				fmt.Println("get floor error")
			}
			floorSlice = append(floorSlice, *floor)
		}
		floors := model.Floors(floorSlice)

		// get datacenter
		dcm := make(map[uint]struct{})
		dcSlice := []model.DataCenter{}
		for _, floor := range floors {
			dcm[floor.DataCenterID] = struct{}{}
		}
		for id := range dcm {
			dc, err := getDataCenter(id)
			if err != nil {
				fmt.Println("get floor error")
			}
			dcSlice = append(dcSlice, *dc)
		}
		dcs := model.DataCenters(dcSlice)

		outputModel := model.Racks{}
		for _, r := range *data {
			row, _ := rows.Take(r.RowID)
			h, _ := halls.Take(row.HallID)
			f, _ := floors.Take(h.FloorID)
			d, _ := dcs.Take(f.DataCenterID)
			f.DataCenter = *d
			h.Floor = *f
			row.Hall = *h
			r.RackRow = *row
			outputModel = append(outputModel, r)
		}

		outputModel.Write(cmd.Flag("output").Value.String())
	}
}

func showUPS(cmd *cobra.Command, args []string) {
	dcName := cmd.Flag("dc").Value.String()
	if len(args) > 0 {
		// show specified row pdu
		url := Conf.APIServer.URL + "/datacenter/ups?name=" + args[0]
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := new(model.UPSs)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("parse response error")
			return
		}
		if len(*data) > 1 {
			fmt.Println("multiple ups found")
		}
		for _, u := range *data {
			dcID := strconv.Itoa(int(u.DataCenterID))
			url = Conf.APIServer.URL + "/datacenter/" + dcID
			body, err := sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if err := json.Unmarshal(body, &u.DataCenter); err != nil {
				fmt.Println(string(body))
				fmt.Println("parse response error")
				return
			}
			u.Write(cmd.Flag("output").Value.String())
		}
	} else {
		// show list of racks
		url := Conf.APIServer.URL + "/datacenter/ups"
		if dcName != "" {
			url = url + "?dc=" + dcName
		}
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := new(model.UPSs)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("parse response error")
			return
		}
		// get dc
		body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter?name="+dcName, []byte{})
		if err != nil {
			fmt.Println(err.Error())
		}
		dcs := new(model.DataCenters)
		if err := json.Unmarshal(body, dcs); err != nil {
			fmt.Println("parse response error")
			return
		}
		outputModel := model.UPSs{}
		for _, u := range *data {
			dc, err := dcs.Take(u.DataCenterID)
			if err != nil {
			}
			u.DataCenter = *dc
			outputModel = append(outputModel, u)
		}

		outputModel.Write(cmd.Flag("output").Value.String())
	}
}

func showRowPDU(cmd *cobra.Command, args []string) {
	dcName := cmd.Flag("dc").Value.String()
	upsName := cmd.Flag("ups").Value.String()
	if len(args) > 0 {
		// show specified row pdu
		url := Conf.APIServer.URL + "/datacenter/row-pdu?name=" + args[0]
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := new(model.RowPDUs)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("parse response error")
			return
		}
		if len(*data) > 1 {
			fmt.Println("multiple row-pdu found")
		}
		for _, p := range *data {
			pUPSID := strconv.Itoa(int(p.PrimaryUPSID))
			url = Conf.APIServer.URL + "/datacenter/ups/" + pUPSID
			body, err := sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if err := json.Unmarshal(body, &p.PrimaryUPS); err != nil {
				fmt.Println(string(body))
				fmt.Println("parse response error")
				return
			}
			if p.SecondaryUPSID != 0 {
				sPDUID := strconv.Itoa(int(p.SecondaryUPSID))
				url = Conf.APIServer.URL + "/datacenter/ups/" + sPDUID
				body, err := sendRequest("GET", url, []byte{})
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				if err := json.Unmarshal(body, &p.SecondaryUPS); err != nil {
					fmt.Println("parse response error")
					return
				}
			}
			p.Write(cmd.Flag("output").Value.String())
		}
	} else {
		// show list of pdus
		url := Conf.APIServer.URL + "/datacenter/row-pdu"
		if dcName != "" {
			url = url + "?dc=" + dcName
			if upsName != "" {
				url = url + "&ups=" + upsName
			}
		}
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := new(model.RowPDUs)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("parse response error")
			return
		}
		// get ups
		body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/ups?dc="+dcName, []byte{})
		if err != nil {
			fmt.Println(err.Error())
		}
		ups := new(model.UPSs)
		if err := json.Unmarshal(body, ups); err != nil {
			fmt.Println("parse response error")
			return
		}
		outputModel := model.RowPDUs{}
		for _, p := range *data {
			pups, err := ups.Take(p.PrimaryUPSID)
			if err != nil {
			}
			p.PrimaryUPS = *pups
			sups, err := ups.Take(p.SecondaryUPSID)
			if err != nil {
			}
			p.SecondaryUPS = *sups
			outputModel = append(outputModel, p)
		}

		outputModel.Write(cmd.Flag("output").Value.String())
	}
}

func showRackPDU(cmd *cobra.Command, args []string) {
	dcName := cmd.Flag("dc").Value.String()
	upsName := cmd.Flag("ups").Value.String()
	pduName := cmd.Flag("pdu").Value.String()
	if len(args) > 0 {
		// show specified rack pdu
		url := Conf.APIServer.URL + "/datacenter/rack-pdu?name=" + args[0]
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := new(model.RackPDUs)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("parse response error")
			return
		}
		if len(*data) > 1 {
			fmt.Println("multiple rack-pdu found")
		}
		for _, p := range *data {
			pPDUID := strconv.Itoa(int(p.PrimaryPDUID))
			url = Conf.APIServer.URL + "/datacenter/row-pdu/" + pPDUID
			body, err := sendRequest("GET", url, []byte{})
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if err := json.Unmarshal(body, &p.PrimaryPDU); err != nil {
				fmt.Println("parse response error")
				return
			}
			if p.SecondaryPDUID != 0 {
				sPDUID := strconv.Itoa(int(p.SecondaryPDUID))
				url = Conf.APIServer.URL + "/datacenter/row-pdu/" + sPDUID
				body, err := sendRequest("GET", url, []byte{})
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				if err := json.Unmarshal(body, &p.SecondaryPDU); err != nil {
					fmt.Println("parse response error")
					return
				}
			} else {
				p.SecondaryPDU = nil
			}
			p.Write(cmd.Flag("output").Value.String())
		}
	} else {
		// show list of pdus
		url := Conf.APIServer.URL + "/datacenter/rack-pdu?"
		if dcName != "" {
			url = url + "&dc=" + dcName
		}
		if upsName != "" {
			url = url + "&ups=" + upsName
		}
		if pduName != "" {
			url = url + "&row-pdu=" + pduName
		}
		body, err := sendRequest("GET", url, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := new(model.RackPDUs)
		if err := json.Unmarshal(body, data); err != nil {
			fmt.Println("parse response error")
			return
		}
		// get row pdu
		rowPduURL := Conf.APIServer.URL + "/datacenter/row-pdu"
		if dcName != "" {
			rowPduURL = rowPduURL + "?dc=" + dcName
		}
		body, err = sendRequest("GET", rowPduURL, []byte{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		dcPDUs := new(model.RowPDUs)
		if err := json.Unmarshal(body, dcPDUs); err != nil {
			fmt.Println("parse response error")
			return
		}
		outputModel := model.RackPDUs{}
		for _, p := range *data {
			pdcpdu, err := dcPDUs.Take(p.PrimaryPDUID)
			if err != nil {
			}
			p.PrimaryPDU = *pdcpdu
			sdcpdu, err := dcPDUs.Take(p.SecondaryPDUID)
			if err != nil {
			}
			if sdcpdu.ID != 0 {
				p.SecondaryPDU = sdcpdu
			} else {
				p.SecondaryPDU = nil
			}
			host, err := getHostByName(p.Name)
			if err != nil {
			}
			rack, _ := getRack(host.RackID)
			loadRackLocation(rack)
			host.Rack = *rack
			p.Host = host
			outputModel = append(outputModel, p)
		}
		outputModel.Write(cmd.Flag("output").Value.String())
	}
}

func createDataCenter(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/datacenter"
	reqModel := model.DataCenter{Name: args[0], Address: cmd.Flag("address").Value.String()}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", url, reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		fmt.Println(reqErr)
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func createDataCenterFloor(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/datacenter"
	// get datacenter
	dcName := cmd.Flag("dc").Value.String()
	body, err := sendRequest("GET", url+"/name/"+dcName, []byte{})
	if err != nil {
		return err
	}
	dc := new(model.DataCenter)
	if err = json.Unmarshal(body, dc); err != nil {
		return fmt.Errorf("response parse error")
	}
	// prepare request floor model
	reqModel := model.Floor{Name: args[0], DataCenterID: dc.ID}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", url+"/floor", reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		fmt.Println(reqErr)
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func createDataCenterHall(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/datacenter"
	// get datacenter
	dcName := cmd.Flag("dc").Value.String()
	body, err := sendRequest("GET", url+"/name/"+dcName, []byte{})
	if err != nil {
		return err
	}
	dc := new(model.DataCenter)
	if err = json.Unmarshal(body, dc); err != nil {
		return fmt.Errorf("response parse error")
	}
	// get datacenter floor
	dcID := strconv.Itoa(int(dc.ID))
	floorName := cmd.Flag("floor").Value.String()
	body, err = sendRequest("GET", url+"/"+dcID+"/floor", []byte{})
	if err != nil {
		return err
	}
	floors := new(model.Floors)
	if err = json.Unmarshal(body, floors); err != nil {
		return fmt.Errorf("response parse error")
	}
	var floorID uint
	for _, floor := range *floors {
		if floor.Name == floorName {
			floorID = floor.ID
			break
		}
	}
	if floorID == 0 {
		return fmt.Errorf("floor not found")
	}
	// hall type
	/*
		hallType := cmd.Flag("type").Value.String()
		if (hallType != "generic") && (hallType != "network") {
			return fmt.Errorf("type must be 'network' or 'generic'")
		}
	*/

	// prepare request hall model
	reqModel := model.Hall{Name: args[0], FloorID: floorID}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", url+"/hall", reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		fmt.Println(reqErr)
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func createRackRow(cmd *cobra.Command, args []string) error {
	// get data hall
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	url := Conf.APIServer.URL + "/datacenter/hall?dc=" + dcName + "&floor=" + floorName + "&name=" + hallName
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	halls := new(model.Halls)
	if err = json.Unmarshal(body, halls); err != nil {
		return fmt.Errorf("response parse error")
	}
	if len(*halls) == 0 {
		return fmt.Errorf("hall not found")
	} else if len(*halls) > 1 {
		return fmt.Errorf("multiple hall found")
	}
	hall := model.Hall{}
	for _, h := range *halls {
		hall = h
	}
	// prepare request hall model
	reqModel := model.RackRow{Name: args[0], HallID: hall.ID}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", Conf.APIServer.URL+"/datacenter/row", reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		fmt.Println(reqErr)
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func createRack(cmd *cobra.Command, args []string) error {
	// get data hall
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	rowName := cmd.Flag("row").Value.String()
	url := Conf.APIServer.URL + "/datacenter/row?dc=" + dcName + "&floor=" + floorName + "&hall=" + hallName + "&name=" + rowName
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	rows := new(model.RackRows)
	if err = json.Unmarshal(body, rows); err != nil {
		return fmt.Errorf("response parse error")
	}
	if len(*rows) == 0 {
		return fmt.Errorf("row not found")
	} else if len(*rows) > 1 {
		return fmt.Errorf("multiple row found")
	}
	row := model.RackRow{}
	for _, r := range *rows {
		row = r
		break
	}
	// prepare request rack model
	reqModel := model.Rack{Name: args[0], RowID: row.ID}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", Conf.APIServer.URL+"/datacenter/rack", reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		fmt.Println(reqErr)
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func createUPS(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/datacenter"
	// get datacenter
	dcName := cmd.Flag("dc").Value.String()
	body, err := sendRequest("GET", url+"/name/"+dcName, []byte{})
	if err != nil {
		return err
	}
	dc := new(model.DataCenter)
	if err = json.Unmarshal(body, dc); err != nil {
		return fmt.Errorf("response parse error")
	}
	// prepare request floor model
	reqModel := model.UPS{Name: args[0], DataCenterID: dc.ID}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", url+"/ups", reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		fmt.Println(reqErr)
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func createPDU(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/datacenter"
	pUPSName := cmd.Flag("primary").Value.String()
	sUPSName := cmd.Flag("secondary").Value.String()
	dcName := cmd.Flag("dc").Value.String()

	// get primary ups
	body, err := sendRequest("GET", url+"/ups?dc="+dcName+"&name="+pUPSName, []byte{})
	if err != nil {
		return err
	}
	upss := new(model.UPSs)
	if err = json.Unmarshal(body, upss); err != nil {
		return fmt.Errorf("response parse error")
	}
	pUPS := model.UPS{}
	for _, u := range *upss {
		pUPS = u
		break
	}

	// get secondary ups
	sUPS := model.UPS{}
	if sUPSName != "" {
		body, err := sendRequest("GET", url+"/ups?dc="+dcName+"&name="+sUPSName, []byte{})
		if err != nil {
			return err
		}
		upss := new(model.UPSs)
		if err = json.Unmarshal(body, upss); err != nil {
			return fmt.Errorf("response parse error")
		}
		for _, u := range *upss {
			sUPS = u
			break
		}
	}

	// prepare request floor model
	reqModel := model.RowPDU{Name: args[0], PrimaryUPSID: pUPS.ID}
	if sUPS.ID != 0 {
		reqModel.SecondaryUPSID = sUPS.ID
	}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", url+"/row-pdu", reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		fmt.Println(reqErr)
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func createRackPDU(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/datacenter"
	pPDUName := cmd.Flag("primary").Value.String()
	sPDUName := cmd.Flag("secondary").Value.String()
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	rowName := cmd.Flag("row").Value.String()
	rackName := cmd.Flag("rack").Value.String()
	groupName := cmd.Flag("group").Value.String()

	// get rack
	rack := new(model.Rack)
	racks, err := getRacks(dcName, floorName, hallName, rowName, rackName)
	if err != nil {
		return err
	}
	for _, r := range *racks {
		rack = &r
		break
	}

	// get primary pdu
	body, err := sendRequest("GET", url+"/row-pdu?dc="+dcName+"&name="+pPDUName, []byte{})
	if err != nil {
		return err
	}
	pdus := new(model.RowPDUs)
	if err = json.Unmarshal(body, pdus); err != nil {
		return fmt.Errorf("response parse error")
	}
	pPDU := model.RowPDU{}
	for _, p := range *pdus {
		pPDU = p
		break
	}
	if pPDU.ID == 0 {
		return fmt.Errorf("primary pdu not found")
	}

	// get secondary pdu
	sPDU := model.RowPDU{}
	if sPDUName != "" {
		body, err := sendRequest("GET", url+"/row-pdu?dc="+dcName+"&name="+sPDUName, []byte{})
		if err != nil {
			return err
		}
		pdus := new(model.RowPDUs)
		if err = json.Unmarshal(body, pdus); err != nil {
			return fmt.Errorf("response parse error")
		}
		for _, p := range *pdus {
			sPDU = p
			break
		}
		if sPDU.ID == 0 {
			return fmt.Errorf("secondary pdu not found")
		}
	}

	// get host group
	group := new(model.HostGroup)
	groups, err := getHostGroups(map[string]string{"name": groupName})
	if err != nil {
		return err
	}
	for _, g := range *groups {
		group = &g
		break
	}

	// prepare request rack-pdu model
	reqModel := model.RackPDU{Name: args[0], PrimaryPDUID: pPDU.ID}
	if sPDU.ID != 0 {
		reqModel.SecondaryPDUID = sPDU.ID
	}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", url+"/rack-pdu", reqJSON)
	if reqErr != nil {
		return reqErr
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	fmt.Println(resMsg.Message)

	// create host
	reqHost := model.Host{Name: args[0], RackID: rack.ID, GroupID: group.ID}
	reqJSON, _ = json.Marshal(reqHost)
	body, err = sendRequest("POST", Conf.APIServer.URL+"/host", reqJSON)
	if err != nil {
		fmt.Println("create host error")
	}
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}

	return nil
}

func updateDataCenter(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/datacenter"
	url = url + "?name=" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	dcs := new(model.DataCenters)
	dc := new(model.DataCenter)
	if err = json.Unmarshal(body, dcs); err != nil {
		return fmt.Errorf("response parse error: " + err.Error())
	}
	for _, d := range *dcs {
		dc = &d
		break
	}
	url = Conf.APIServer.URL + "/datacenter/" + strconv.Itoa(int(dc.ID))
	reqModel := model.DataCenter{Address: cmd.Flag("address").Value.String()}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("PUT", url, reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		fmt.Println(reqErr)
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func updateDataCenterFloor(cmd *cobra.Command, args []string) error {
	floorName := cmd.Flag("name").Value.String()
	dcName := cmd.Flag("dc").Value.String()
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/name/"+dcName, []byte{})
	if err != nil {
		return err
	}
	dc := new(model.DataCenter)
	if err = json.Unmarshal(body, dc); err != nil {
		return fmt.Errorf("response parse error: " + err.Error())
	}
	dcID := strconv.Itoa(int(dc.ID))
	body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
	if err != nil {
		return err
	}
	data := new(model.Floors)
	if err := json.Unmarshal(body, data); err != nil {
		return fmt.Errorf("response parse error:" + err.Error())
	}
	floor := model.Floor{}
	for _, flr := range *data {
		if flr.Name == args[0] {
			floor = flr
		}
	}
	if floor.ID == 0 {
		return fmt.Errorf("floor not found")
	}
	if floorName == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	floor.Name = floorName
	reqJSON, _ := json.Marshal(floor)
	floorID := strconv.Itoa(int(floor.ID))
	url := Conf.APIServer.URL + "/datacenter/floor/" + floorID
	body, err = sendRequest("PUT", url, reqJSON)
	if err != nil {
		return err
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	fmt.Println(resMsg.Message)

	return nil
}

func updateDataCenterHall(cmd *cobra.Command, args []string) error {
	hallName := cmd.Flag("name").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	dcName := cmd.Flag("dc").Value.String()
	if hallName == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/name/"+dcName, []byte{})
	if err != nil {
		return err
	}

	dc := new(model.DataCenter)
	if err = json.Unmarshal(body, dc); err != nil {
		return fmt.Errorf("response parse error: " + err.Error())
	}
	dcID := strconv.Itoa(int(dc.ID))
	body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
	if err != nil {
		return err
	}
	data := new(model.Floors)
	if err := json.Unmarshal(body, data); err != nil {
		return fmt.Errorf("response parse error:" + err.Error())
	}
	floor := model.Floor{}
	for _, flr := range *data {
		if flr.Name == floorName {
			floor = flr
			break
		}
	}
	if floor.ID == 0 {
		return fmt.Errorf("floor not found")
	}
	floorID := strconv.Itoa(int(floor.ID))
	url := Conf.APIServer.URL + "/datacenter/floor/" + floorID
	body, err = sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	f := new(model.Floor)
	if err := json.Unmarshal(body, f); err != nil {
		return fmt.Errorf("response parse error: " + err.Error())
	}
	hall := model.Hall{}
	for _, h := range f.Halls {
		if h.Name == args[0] {
			hall = h
			break
		}
	}
	if hall.ID == 0 {
		return fmt.Errorf("hall not found")
	}
	hall.Name = hallName
	reqJSON, _ := json.Marshal(hall)
	hallID := strconv.Itoa(int(hall.ID))
	url = Conf.APIServer.URL + "/datacenter/hall/" + hallID
	body, err = sendRequest("PUT", url, reqJSON)
	if err != nil {
		return err
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	fmt.Println(resMsg.Message)

	return nil
}

func updateRackRow(cmd *cobra.Command, args []string) error {
	rowName := cmd.Flag("name").Value.String()
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	if rowName == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	url := Conf.APIServer.URL + "/datacenter/row?dc=" + dcName + "&floor=" + floorName + "&hall=" + hallName + "&name=" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	rows := new([]model.RackRow)
	if err := json.Unmarshal(body, rows); err != nil {
		return fmt.Errorf("response parse error" + err.Error())
	}
	if len(*rows) > 1 {
		return fmt.Errorf("multiple row found")
	}
	row := new(model.RackRow)
	for _, r := range *rows {
		row = &r
		break
	}
	row.Name = rowName
	reqJSON, _ := json.Marshal(row)
	rowID := strconv.Itoa(int(row.ID))
	url = Conf.APIServer.URL + "/datacenter/row/" + rowID
	body, reqErr := sendRequest("PUT", url, reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func updateRack(cmd *cobra.Command, args []string) error {
	rackName := cmd.Flag("name").Value.String()
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	rowName := cmd.Flag("row").Value.String()
	if rackName == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	url := Conf.APIServer.URL + "/datacenter/rack?dc=" + dcName + "&floor=" + floorName + "&hall=" + hallName + "&row=" + rowName + "&name=" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	racks := new(model.Racks)
	if err := json.Unmarshal(body, racks); err != nil {
		return fmt.Errorf("response parse error" + err.Error())
	}
	if len(*racks) > 1 {
		return fmt.Errorf("multiple rack found")
	}
	rack := new(model.Rack)
	for _, r := range *racks {
		rack = &r
		break
	}
	rack.Name = rackName
	reqJSON, _ := json.Marshal(rack)
	rackID := strconv.Itoa(int(rack.ID))
	url = Conf.APIServer.URL + "/datacenter/rack/" + rackID
	body, reqErr := sendRequest("PUT", url, reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func updateUPS(cmd *cobra.Command, args []string) error {
	upsName := cmd.Flag("name").Value.String()
	dcName := cmd.Flag("dc").Value.String()
	if upsName == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	url := Conf.APIServer.URL + "/datacenter/ups?dc=" + dcName + "&name=" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	upss := new(model.UPSs)
	if err := json.Unmarshal(body, upss); err != nil {
		return fmt.Errorf("response parse error" + err.Error())
	}
	if len(*upss) > 1 {
		return fmt.Errorf("multiple ups found")
	}
	ups := new(model.UPS)
	for _, r := range *upss {
		ups = &r
		break
	}
	ups.Name = upsName
	reqJSON, _ := json.Marshal(ups)
	upsID := strconv.Itoa(int(ups.ID))
	url = Conf.APIServer.URL + "/datacenter/ups/" + upsID
	body, reqErr := sendRequest("PUT", url, reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func updatePDU(cmd *cobra.Command, args []string) error {
	pduName := cmd.Flag("name").Value.String()
	dcName := cmd.Flag("dc").Value.String()
	if pduName == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	url := Conf.APIServer.URL + "/datacenter/row-pdu?dc=" + dcName + "&name=" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	pdus := new(model.RowPDUs)
	if err := json.Unmarshal(body, pdus); err != nil {
		return fmt.Errorf("response parse error" + err.Error())
	}
	if len(*pdus) > 1 {
		return fmt.Errorf("multiple pdu found")
	}
	pdu := new(model.RowPDU)
	for _, p := range *pdus {
		pdu = &p
		break
	}
	pdu.Name = pduName
	reqJSON, _ := json.Marshal(pdu)
	pduID := strconv.Itoa(int(pdu.ID))
	url = Conf.APIServer.URL + "/datacenter/row-pdu/" + pduID
	body, reqErr := sendRequest("PUT", url, reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func updateRackPDU(cmd *cobra.Command, args []string) error {
	pduName := cmd.Flag("name").Value.String()
	dcName := cmd.Flag("dc").Value.String()
	if pduName == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	url := Conf.APIServer.URL + "/datacenter/rack-pdu?dc=" + dcName + "&name=" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	pdus := new(model.RackPDUs)
	if err := json.Unmarshal(body, pdus); err != nil {
		return fmt.Errorf("response parse error" + err.Error())
	}
	if len(*pdus) > 1 {
		return fmt.Errorf("multiple pdu found")
	}
	pdu := new(model.RackPDU)
	for _, p := range *pdus {
		pdu = &p
		break
	}
	pdu.Name = pduName
	reqJSON, _ := json.Marshal(pdu)
	pduID := strconv.Itoa(int(pdu.ID))
	url = Conf.APIServer.URL + "/datacenter/rack-pdu/" + pduID
	body, reqErr := sendRequest("PUT", url, reqJSON)
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteDataCenter(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/datacenter/" + args[0]
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		fmt.Println(reqErr)
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteDataCenterFloor(cmd *cobra.Command, args []string) error {
	dcName := cmd.Flag("dc").Value.String()
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/name/"+dcName, []byte{})
	if err != nil {
		return err
	}
	dc := new(model.DataCenter)
	if err = json.Unmarshal(body, dc); err != nil {
		return fmt.Errorf("response parse error: " + err.Error())
	}
	dcID := strconv.Itoa(int(dc.ID))
	body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
	if err != nil {
		return err
	}
	data := new(model.Floors)
	if err := json.Unmarshal(body, data); err != nil {
		return fmt.Errorf("response parse error:" + err.Error())
	}
	floor := model.Floor{}
	for _, flr := range *data {
		if flr.Name == args[0] {
			floor = flr
		}
	}
	if floor.ID == 0 {
		return fmt.Errorf("floor not found")
	}
	floorID := strconv.Itoa(int(floor.ID))
	url := Conf.APIServer.URL + "/datacenter/floor/" + floorID
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteDataCenterHall(cmd *cobra.Command, args []string) error {
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/name/"+dcName, []byte{})
	if err != nil {
		return err
	}
	dc := new(model.DataCenter)
	if err = json.Unmarshal(body, dc); err != nil {
		return fmt.Errorf("response parse error: " + err.Error())
	}
	dcID := strconv.Itoa(int(dc.ID))
	body, err = sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+dcID+"/floor", []byte{})
	if err != nil {
		return err
	}
	data := new(model.Floors)
	if err := json.Unmarshal(body, data); err != nil {
		return fmt.Errorf("response parse error:" + err.Error())
	}
	floor := model.Floor{}
	for _, flr := range *data {
		if flr.Name == floorName {
			floor = flr
		}
	}
	if floor.ID == 0 {
		return fmt.Errorf("floor not found")
	}
	floorID := strconv.Itoa(int(floor.ID))
	url := Conf.APIServer.URL + "/datacenter/floor/" + floorID

	f := new(model.Floor)
	body, err = sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, f); err != nil {
		return err
	}

	hall := model.Hall{}
	for _, h := range f.Halls {
		if h.Name == args[0] {
			hall = h
			break
		}
	}
	if hall.ID == 0 {
		return fmt.Errorf("hall not found")
	}
	hallID := strconv.Itoa(int(hall.ID))
	url = Conf.APIServer.URL + "/datacenter/hall/" + hallID
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteRackRow(cmd *cobra.Command, args []string) error {
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()

	url := Conf.APIServer.URL + "/datacenter/row?dc=" + dcName + "&floor=" + floorName + "&hall=" + hallName + "&name=" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	rows := new([]model.RackRow)
	if err := json.Unmarshal(body, rows); err != nil {
		return fmt.Errorf("response parse error" + err.Error())
	}
	if len(*rows) > 1 {
		return fmt.Errorf("multiple row found")
	}
	row := new(model.RackRow)
	for _, r := range *rows {
		row = &r
		break
	}
	rowID := strconv.Itoa(int(row.ID))
	url = Conf.APIServer.URL + "/datacenter/row/" + rowID
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteRack(cmd *cobra.Command, args []string) error {
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	hallName := cmd.Flag("hall").Value.String()
	rowName := cmd.Flag("row").Value.String()

	url := Conf.APIServer.URL + "/datacenter/rack?dc=" + dcName + "&floor=" + floorName + "&hall=" + hallName + "&row=" + rowName + "&name=" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	racks := new(model.Racks)
	if err := json.Unmarshal(body, racks); err != nil {
		return fmt.Errorf("response parse error" + err.Error())
	}
	if len(*racks) > 1 {
		return fmt.Errorf("multiple rack found")
	}
	rack := new(model.Rack)
	for _, r := range *racks {
		rack = &r
		break
	}
	rackID := strconv.Itoa(int(rack.ID))
	url = Conf.APIServer.URL + "/datacenter/rack/" + rackID
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteUPS(cmd *cobra.Command, args []string) error {
	dcName := cmd.Flag("dc").Value.String()
	url := Conf.APIServer.URL + "/datacenter/ups?dc=" + dcName + "&name=" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	upss := new(model.UPSs)
	if err := json.Unmarshal(body, upss); err != nil {
		return fmt.Errorf("response parse error" + err.Error())
	}
	if len(*upss) > 1 {
		return fmt.Errorf("multiple ups found")
	}
	ups := new(model.UPS)
	for _, u := range *upss {
		ups = &u
		break
	}
	upsID := strconv.Itoa(int(ups.ID))
	url = Conf.APIServer.URL + "/datacenter/ups/" + upsID
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deletePDU(cmd *cobra.Command, args []string) error {
	dcName := cmd.Flag("dc").Value.String()
	url := Conf.APIServer.URL + "/datacenter/row-pdu?dc=" + dcName + "&name=" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	pdus := new(model.RowPDUs)
	if err := json.Unmarshal(body, pdus); err != nil {
		return fmt.Errorf("response parse error" + err.Error())
	}
	if len(*pdus) > 1 {
		return fmt.Errorf("multiple pdu found")
	}
	pdu := new(model.RowPDU)
	for _, p := range *pdus {
		pdu = &p
		break
	}
	pduID := strconv.Itoa(int(pdu.ID))
	url = Conf.APIServer.URL + "/datacenter/row-pdu/" + pduID
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

func deleteRackPDU(cmd *cobra.Command, args []string) error {
	dcName := cmd.Flag("dc").Value.String()
	url := Conf.APIServer.URL + "/datacenter/rack-pdu?dc=" + dcName + "&name=" + args[0]
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return err
	}
	pdus := new(model.RackPDUs)
	if err := json.Unmarshal(body, pdus); err != nil {
		return fmt.Errorf("response parse error" + err.Error())
	}
	if len(*pdus) > 1 {
		return fmt.Errorf("multiple rack pdu found")
	}
	pdu := new(model.RackPDU)
	for _, p := range *pdus {
		pdu = &p
		break
	}
	pduID := strconv.Itoa(int(pdu.ID))
	url = Conf.APIServer.URL + "/datacenter/rack-pdu/" + pduID
	body, reqErr := sendRequest("DELETE", url, []byte{})
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	if reqErr != nil {
		return reqErr
	}
	fmt.Println(resMsg.Message)

	return nil
}

//
func getDataCenter(id uint) (*model.DataCenter, error) {
	dc := new(model.DataCenter)
	idStr := strconv.Itoa(int(id))
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/"+idStr, []byte{})
	if err != nil {
		return dc, err
	}
	if err := json.Unmarshal(body, dc); err != nil {
		return dc, fmt.Errorf("response parse error")
	}

	return dc, nil
}

func getFloor(id uint) (*model.Floor, error) {
	floor := new(model.Floor)
	idStr := strconv.Itoa(int(id))
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/floor/"+idStr, []byte{})
	if err != nil {
		return floor, err
	}
	if err := json.Unmarshal(body, floor); err != nil {
		return floor, fmt.Errorf("response parse error")
	}

	return floor, nil
}

func getHall(id uint) (*model.Hall, error) {
	hall := new(model.Hall)
	idStr := strconv.Itoa(int(id))
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/hall/"+idStr, []byte{})
	if err != nil {
		return hall, err
	}
	if err := json.Unmarshal(body, hall); err != nil {
		return hall, fmt.Errorf("response parse error")
	}

	return hall, nil
}

func getRow(id uint) (*model.RackRow, error) {
	row := new(model.RackRow)
	idStr := strconv.Itoa(int(id))
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/row/"+idStr, []byte{})
	if err != nil {
		return row, err
	}
	if err := json.Unmarshal(body, row); err != nil {
		return row, fmt.Errorf("response parse error")
	}

	return row, nil
}

func getRack(id uint) (*model.Rack, error) {
	rack := new(model.Rack)
	idStr := strconv.Itoa(int(id))
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/rack/"+idStr, []byte{})
	if err != nil {
		return rack, err
	}
	if err := json.Unmarshal(body, rack); err != nil {
		return rack, fmt.Errorf("response parse error")
	}

	return rack, nil
}

func getRacks(dc, floor, hall, row, name string) (*model.Racks, error) {
	racks := new(model.Racks)
	url := Conf.APIServer.URL + "/datacenter/rack"
	url = url + "?dc=" + dc + "&floor=" + floor + "&hall=" + hall + "&row=" + row + "&name=" + name
	body, err := sendRequest("GET", url, []byte{})
	if err != nil {
		return &model.Racks{}, err
	}
	if err := json.Unmarshal(body, racks); err != nil {
		return &model.Racks{}, fmt.Errorf("response parse error")
	}
	if len(*racks) == 0 {
		return racks, fmt.Errorf("rack not found")
	}

	return racks, nil
}

func loadRackLocation(rack *model.Rack) {
	row, _ := getRow(rack.RowID)
	hall, _ := getHall(row.HallID)
	floor, _ := getFloor(hall.FloorID)
	dc, _ := getDataCenter(floor.DataCenterID)
	floor.DataCenter = *dc
	hall.Floor = *floor
	row.Hall = *hall
	rack.RackRow = *row
}
