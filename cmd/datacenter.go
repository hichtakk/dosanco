package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hichtakk/dosanco/model"
)

func showDataCenter(cmd *cobra.Command, args []string) {
	tree := cmd.Flag("tree").Value.String()
	url := Conf.APIServer.URL + "/datacenter"
	if len(args) > 0 {
		// show specified datacenter
		dcs, _ := getDataCenters(map[string]string{"name": args[0]})
		if len(*dcs) == 0 {
			fmt.Println("datacenter not found")
			return
		}
		dc := new(model.DataCenter)
		for _, d := range *dcs {
			dc = &d
		}
		if tree == "true" {
			floors := new(model.Floors)
			dcFloors, _ := getFloors(map[string]string{"dc": dc.Name})
			for _, dcFloor := range *dcFloors {
				halls := new(model.Halls)
				floorHalls, _ := getHalls(map[string]string{"dc": dc.Name, "floor": dcFloor.Name})
				for _, floorHall := range *floorHalls {
					rows := new(model.RackRows)
					hallRows, _ := getRows(map[string]string{"dc": dc.Name, "floor": dcFloor.Name, "hall": floorHall.Name})
					for _, hallRow := range *hallRows {
						racks := new(model.Racks)
						rowRacks, _ := getRacks(map[string]string{"dc": dc.Name, "floor": dcFloor.Name, "hall": floorHall.Name, "row": hallRow.Name})
						for _, rack := range *rowRacks {
							*racks = append(*racks, rack)
						}
						hallRow.Racks = *racks
						*rows = append(*rows, hallRow)
					}
					floorHall.RackRows = *rows
					*halls = append(*halls, floorHall)
				}
				dcFloor.Halls = *halls
				*floors = append(*floors, dcFloor)
			}
			dc.Floors = floors
			dc.WriteTree(cmd.Flag("output").Value.String())
		} else {
			floors, _ := getFloors(map[string]string{"dc": dc.Name})
			dc.Floors = floors
			dc.Write(cmd.Flag("output").Value.String())
		}
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
		output := new(model.DataCenters)
		if tree == "true" {
			for _, dc := range *data {
				floors := new(model.Floors)
				dcFloors, _ := getFloors(map[string]string{"dc": dc.Name})
				for _, dcFloor := range *dcFloors {
					halls := new(model.Halls)
					floorHalls, _ := getHalls(map[string]string{"dc": dc.Name, "floor": dcFloor.Name})
					for _, floorHall := range *floorHalls {
						rows := new(model.RackRows)
						hallRows, _ := getRows(map[string]string{"dc": dc.Name, "floor": dcFloor.Name, "hall": floorHall.Name})
						for _, hallRow := range *hallRows {
							racks := new(model.Racks)
							rowRacks, _ := getRacks(map[string]string{"dc": dc.Name, "floor": dcFloor.Name, "hall": floorHall.Name, "row": hallRow.Name})
							for _, rack := range *rowRacks {
								*racks = append(*racks, rack)
							}
							hallRow.Racks = *racks
							*rows = append(*rows, hallRow)
						}
						floorHall.RackRows = *rows
						*halls = append(*halls, floorHall)
					}
					dcFloor.Halls = *halls
					*floors = append(*floors, dcFloor)
				}
				dc.Floors = floors
				*output = append(*output, dc)
			}
			output.WriteTree(cmd.Flag("output").Value.String())
		} else {
			data.Write(cmd.Flag("output").Value.String())
		}
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
		floor := new(model.Floor)
		floors, _ := getFloors(map[string]string{"dc": dcName, "name": args[0]})
		if len(*floors) == 0 {
			fmt.Println(fmt.Sprintf("floor '%v' not found at '%v'", args[0], dcName))
		}
		for _, f := range *floors {
			floor = &f
			break
		}
		floor.Write(cmd.Flag("output").Value.String())
	} else {
		if dcName != "" {
			// show all datacenter floors
			dcs, _ := getDataCenters(map[string]string{"name": dcName})
			if len(*dcs) == 0 {
				fmt.Println("datacenter not found")
				return
			}
			dc := new(model.DataCenter)
			for _, d := range *dcs {
				dc = &d
			}
			floors, _ := getFloors(map[string]string{"dc": dc.Name})
			if len(*floors) == 0 {
				fmt.Println("floor not found")
				return
			}
			output := model.Floors{}
			for _, f := range *floors {
				f.DataCenter = dc
				output = append(output, f)
			}
			output.Write(cmd.Flag("output").Value.String())
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
				f.DataCenter = dc
				outputModel = append(outputModel, f)
			}

			outputModel.Write(cmd.Flag("output").Value.String())
		}
	}
}

func showDataCenterHall(cmd *cobra.Command, args []string) {
	//url := Conf.APIServer.URL + "/datacenter/hall"
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
		hall := new(model.Hall)
		halls, err := getHalls(map[string]string{"dc": dcName, "floor": floorName, "name": args[0]})
		if err != nil {
			fmt.Println(err)
		}
		for _, h := range *halls {
			floor, _ := getFloor(h.FloorID)
			h.Floor = floor
			hall = &h
		}
		hall.Write(cmd.Flag("output").Value.String())
	} else {
		if dcName != "" {
			// show all halls of specified datacenter
			floors, _ := getFloors(map[string]string{"dc": dcName})
			if floorName == "" {
				halls := model.Halls{}
				for _, floor := range *floors {
					flr := floor
					hls, _ := getHalls(map[string]string{"dc": dcName, "floor": floor.Name})
					for _, h := range *hls {
						h.Floor = &flr
						halls = append(halls, h)
					}
				}
				halls.Write(cmd.Flag("output").Value.String())
			} else {
				floor := new(model.Floor)
				for _, f := range *floors {
					if f.Name == floorName {
						floor = &f
						break
					}
				}
				fID := strconv.Itoa(int(floor.ID))
				body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/floor/"+fID, []byte{})
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
					h.Floor = floor
					outputModel = append(outputModel, h)
				}
				outputModel.Write(cmd.Flag("output").Value.String())
			}
		} else {
			// show all datacenter halls
			if floorName != "" {
				// TODO: get halls specifing only floor name
				/*
					floor, _ := getFloors(map[string]string{"name": floorName})
					for _, f := range *floor {
					}
				*/
			} else {
				outputModel := model.Halls{}
				dcs, _ := getDataCenters(map[string]string{})
				for _, dc := range *dcs {
					floors, _ := getFloors(map[string]string{"dc": dc.Name})
					for _, floor := range *floors {
						halls, _ := getHalls(map[string]string{"dc": dc.Name, "floor": floor.Name})
						for _, hall := range *halls {
							flr := floor
							hall.Floor = &flr
							outputModel = append(outputModel, hall)
						}
					}
				}
				outputModel.Write(cmd.Flag("output").Value.String())
			}
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
			f.DataCenter = d
			hall.Floor = f
			r.Hall = hall
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
			f.DataCenter = d
			h.Floor = f
			row.Hall = h
			r.RackRow = row
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
	location := cmd.Flag("location").Value.String()
	if len(args) > 0 {
		pdus, err := getRackPDUs(map[string]string{"name": args[0]})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, p := range *pdus {
			if p.PrimaryPDUID != 0 {
				pPDU, err := getRowPDU(p.PrimaryPDUID)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				p.PrimaryPDU = pPDU
			}
			if p.SecondaryPDUID != 0 {
				sPDU, err := getRowPDU(p.SecondaryPDUID)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				p.SecondaryPDU = sPDU
			}
			hosts, _ := getHosts(map[string]string{"name": p.Name})
			for _, h := range *hosts {
				grp, _ := getHostGroup(h.GroupID)
				h.Group = grp
				rck, _ := getRack(h.RackID)
				loadRackLocation(rck)
				h.Rack = *rck
				allocs, _ := getIPv4Allocations(map[string]string{"name": p.Name})
				for _, alloc := range *allocs {
					nw, _ := getNetwork(alloc.IPv4NetworkID)
					network := *nw
					alloc.IPv4Network = &network
					h.IPv4Allocations = append(h.IPv4Allocations, alloc)
				}
				p.Host = &h
				break
			}
			p.Write(cmd.Flag("output").Value.String())
		}
	} else {
		// show list of pdus
		if location != "" {
			query := map[string]string{}
			query["location"] = url.QueryEscape(location)
			query["type"] = "rack-pdu"
			hosts, err := getHosts(query)
			if err != nil {
				fmt.Println(err)
				return
			}
			output := new(model.RackPDUs)
			for _, h := range *hosts {
				rack, _ := getRack(h.RackID)
				loadRackLocation(rack)
				h.Rack = *rack
				pdus, _ := getRackPDUs(map[string]string{"name": h.Name})
				for _, p := range *pdus {
					pdu := p
					pduHost := h
					pdu.Host = &pduHost
					rowPDU, _ := getRowPDU(pdu.PrimaryPDUID)
					pdu.PrimaryPDU = rowPDU
					if pdu.SecondaryPDUID != 0 {
						srowPDU, _ := getRowPDU(pdu.SecondaryPDUID)
						pdu.SecondaryPDU = srowPDU
					}
					*output = append(*output, pdu)
				}
			}
			output.Write(cmd.Flag("output").Value.String())
		} else {
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
			dcPDUs, err := getRowPDUs(map[string]string{"dc": dcName})
			if err != nil {
				fmt.Printf("getting row-pdu error\n")
				return
			}
			outputModel := model.RackPDUs{}
			for _, p := range *data {
				if p.PrimaryPDUID != 0 {
					pdcpdu, _ := dcPDUs.Take(p.PrimaryPDUID)
					p.PrimaryPDU = pdcpdu
				} else {
					p.PrimaryPDU = nil
				}
				if p.SecondaryPDUID != 0 {
					sdcpdu, _ := dcPDUs.Take(p.SecondaryPDUID)
					p.SecondaryPDU = sdcpdu
				} else {
					p.SecondaryPDU = nil
				}
				host := new(model.Host)
				hosts, _ := getHosts(map[string]string{"name": p.Name})
				for _, h := range *hosts {
					host = &h
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
	dc := new(model.DataCenter)
	dcs, _ := getDataCenters(map[string]string{"name": dcName})
	for _, d := range *dcs {
		dc = &d
		break
	}
	// prepare request floor model
	reqModel := model.Floor{Name: args[0], DataCenterID: dc.ID}
	reqJSON, err := json.Marshal(reqModel)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", reqModel)
	}
	body, reqErr := sendRequest("POST", url+"/floor", reqJSON)
	if reqErr != nil {
		return reqErr
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
	}
	fmt.Println(resMsg.Message)

	return nil
}

func createDataCenterHall(cmd *cobra.Command, args []string) error {
	url := Conf.APIServer.URL + "/datacenter"
	// get datacenter
	dcName := cmd.Flag("dc").Value.String()
	floorName := cmd.Flag("floor").Value.String()
	floor := new(model.Floor)
	floors, _ := getFloors(map[string]string{"dc": dcName, "name": floorName})
	if len(*floors) == 0 {
		return fmt.Errorf("floor not found")
	} else if len(*floors) > 1 {
		return fmt.Errorf("multiple floors are found")
	}
	for _, f := range *floors {
		floor = &f
		break
	}
	// prepare request hall model
	reqModel := model.Hall{Name: args[0], FloorID: floor.ID}
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
	description := cmd.Flag("description").Value.String()
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
	reqModel := model.Rack{Name: args[0], RowID: row.ID, Description: description}
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
	dc := new(model.DataCenter)
	dcs, _ := getDataCenters(map[string]string{"name": dcName})
	for _, d := range *dcs {
		dc = &d
		break
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
	location := cmd.Flag("location").Value.String()
	groupName := cmd.Flag("group").Value.String()
	locSlice := strings.Split(location, "/")
	if len(locSlice) != 5 {
		return fmt.Errorf("invalid location format. use '{DC}/{FLOOR}/{HALL}/{ROW}/{RACK}'")
	}
	dcName := locSlice[0]
	floorName := locSlice[1]
	hallName := locSlice[2]
	rowName := locSlice[3]
	rackName := locSlice[4]
	racks, err := getRacks(map[string]string{"dc": dcName, "floor": floorName, "hall": hallName, "row": rowName, "name": rackName})
	if err != nil {
		return fmt.Errorf("rack not found for specified location")
	}
	rack := new(model.Rack)
	for _, r := range *racks {
		rack = &r
		break
	}

	// get primary pdu
	pPDU := model.RowPDU{}
	if pPDUName != "" {
		body, err := sendRequest("GET", url+"/row-pdu?dc="+dcName+"&name="+pPDUName, []byte{})
		if err != nil {
			return err
		}
		pdus := new(model.RowPDUs)
		if err = json.Unmarshal(body, pdus); err != nil {
			return fmt.Errorf("response parse error")
		}
		for _, p := range *pdus {
			pPDU = p
			break
		}
		if pPDU.ID == 0 {
			return fmt.Errorf("primary pdu not found")
		}
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
	reqModel := model.RackPDU{Name: args[0]}
	if pPDU.ID != 0 {
		reqModel.PrimaryPDUID = pPDU.ID
	}
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
	reqHost := model.Host{Name: args[0], RackID: rack.ID, GroupID: group.ID, Type: "rack-pdu"}
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
	if floorName == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	floors, err := getFloors(map[string]string{"dc": dcName, "name": args[0]})
	if err != nil {
		return err
	}
	floor := new(model.Floor)
	for _, f := range *floors {
		floor = &f
		break
	}
	floor.Name = floorName
	reqJSON, _ := json.Marshal(floor)
	floorID := strconv.Itoa(int(floor.ID))
	url := Conf.APIServer.URL + "/datacenter/floor/" + floorID
	body, err := sendRequest("PUT", url, reqJSON)
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
	hall := new(model.Hall)
	halls, err := getHalls(map[string]string{"dc": dcName, "floor": floorName, "name": args[0]})
	if err != nil {
		return err
	}
	if len(*halls) == 0 {
		return fmt.Errorf("hall not found")
	}
	for _, h := range *halls {
		hall = &h
		break
	}
	hall.Name = hallName
	reqJSON, _ := json.Marshal(hall)
	hallID := strconv.Itoa(int(hall.ID))
	url := Conf.APIServer.URL + "/datacenter/hall/" + hallID
	body, err := sendRequest("PUT", url, reqJSON)
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
	description := cmd.Flag("description").Value.String()
	if rackName == "-" && description == "-" {
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
	if rackName != "-" {
		rack.Name = rackName
	}
	if description != "-" {
		rack.Description = description
	}
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
	description := cmd.Flag("description").Value.String()
	if upsName == "-" && description == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	ups := new(model.UPS)
	upss, err := getUPSs(map[string]string{"dc": dcName, "name": args[0]})
	if err != nil {
		return err
	}
	if len(*upss) == 0 {
		return fmt.Errorf("ups not found")
	}
	if len(*upss) > 1 {
		return fmt.Errorf("multiple ups are found")
	}
	for _, u := range *upss {
		ups = &u
		break
	}
	if upsName != "-" {
		ups.Name = upsName
	}
	if description != "-" {
		ups.Description = description
	}
	reqJSON, _ := json.Marshal(ups)
	upsID := strconv.Itoa(int(ups.ID))
	url := Conf.APIServer.URL + "/datacenter/ups/" + upsID
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
	description := cmd.Flag("description").Value.String()
	if pduName == "-" && description == "-" {
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
	if pduName != "-" {
		pdu.Name = pduName
	}
	if description != "-" {
		pdu.Description = description
	}
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
	primary := cmd.Flag("primary").Value.String()
	secondary := cmd.Flag("secondary").Value.String()
	description := cmd.Flag("description").Value.String()
	if pduName == "-" && primary == "-" && secondary == "-" && description == "-" {
		return fmt.Errorf("nothing to be updated")
	}
	/*
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
	*/
	pdu := new(model.RackPDU)
	pdus, err := getRackPDUs(map[string]string{"name": args[0], "dc": dcName})
	if err != nil {
		return err
	}
	for _, p := range *pdus {
		pdu = &p
		break
	}

	if pduName != "-" {
		if exist, _ := getHosts(map[string]string{"name": pduName}); len(*exist) > 0 {
			return fmt.Errorf("new name '%v' is already exist", pduName)
		}
		pduHost := new(model.Host)
		pduHosts, _ := getHosts(map[string]string{"name": pdu.Name})
		for _, h := range *pduHosts {
			pduHost = &h
			break
		}
		pduHost.Name = pduName
		hostReqJSON, _ := json.Marshal(pduHost)
		hostID := strconv.Itoa(int(pduHost.ID))
		body, _ := sendRequest("PUT", Conf.APIServer.URL+"/host/"+hostID, hostReqJSON)
		var hostReqMsg responseMessage
		if err := json.Unmarshal(body, &hostReqMsg); err != nil {
			return err
		}
		pdu.Name = pduName
	}
	if primary != "-" {
		rowPDUs, _ := getRowPDUs(map[string]string{"name": primary})
		if len(*rowPDUs) == 0 {
			return fmt.Errorf("row-pdu '%v' not found", primary)
		}
		if len(*rowPDUs) > 1 {
			return fmt.Errorf("multiple row-pdu found")
		}
		for _, rp := range *rowPDUs {
			pdu.PrimaryPDUID = rp.ID
			break
		}
	}
	if secondary != "-" && secondary != "" {
		rowPDUs, _ := getRowPDUs(map[string]string{"name": secondary})
		if len(*rowPDUs) == 0 {
			return fmt.Errorf("row-pdu '%v' not found", primary)
		}
		if len(*rowPDUs) > 1 {
			return fmt.Errorf("multiple row-pdu found")
		}
		for _, rp := range *rowPDUs {
			pdu.SecondaryPDUID = rp.ID
			break
		}
	} else if secondary == "" {
		pdu.SecondaryPDUID = 0
	}
	if description != "-" {
		pdu.Description = description
	}

	reqJSON, _ := json.Marshal(pdu)
	pduID := strconv.Itoa(int(pdu.ID))
	url := Conf.APIServer.URL + "/datacenter/rack-pdu/" + pduID
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
	floor := new(model.Floor)
	floors, _ := getFloors(map[string]string{"dc": dcName, "name": args[0]})
	for _, f := range *floors {
		floor = &f
		break
	}
	if len(*floors) == 0 {
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
	hall := new(model.Hall)
	halls, _ := getHalls(map[string]string{"dc": dcName, "floor": floorName, "name": args[0]})
	for _, h := range *halls {
		hall = &h
		break
	}
	if len(*halls) == 0 {
		return fmt.Errorf("hall not found")
	}
	hallID := strconv.Itoa(int(hall.ID))
	url := Conf.APIServer.URL + "/datacenter/hall/" + hallID
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

	racks, err := getRacks(map[string]string{"dc": dcName, "floor": floorName, "hall": hallName, "row": rowName, "name": args[0]})
	if err != nil {
		return err
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
	url := Conf.APIServer.URL + "/datacenter/rack/" + rackID
	body, reqErr := sendRequest("DELETE", url, []byte{})
	if reqErr != nil {
		return reqErr
	}
	var resMsg responseMessage
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return err
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

func getDataCenters(query map[string]string) (*model.DataCenters, error) {
	dcs := new(model.DataCenters)
	queryString := ""
	for key, val := range query {
		queryString = queryString + "&" + key + "=" + val
	}
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter?"+queryString, []byte{})
	if err != nil {
		return dcs, err
	}
	if err := json.Unmarshal(body, dcs); err != nil {
		return dcs, fmt.Errorf("response parse error")
	}

	return dcs, nil
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

func getFloors(query map[string]string) (*model.Floors, error) {
	floors := new(model.Floors)
	queryString := ""
	for key, val := range query {
		queryString = queryString + "&" + key + "=" + val
	}
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/floor?"+queryString, []byte{})
	if err != nil {
		return floors, err
	}
	if err := json.Unmarshal(body, floors); err != nil {
		return floors, fmt.Errorf("response parse error")
	}
	if len(*floors) == 0 {
		return floors, fmt.Errorf("floor not found")
	}

	return floors, nil
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

func getHalls(query map[string]string) (*model.Halls, error) {
	halls := new(model.Halls)
	queryString := ""
	for key, val := range query {
		queryString = queryString + "&" + key + "=" + val
	}
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/hall?"+queryString, []byte{})
	if err != nil {
		return halls, err
	}
	if err := json.Unmarshal(body, halls); err != nil {
		return halls, fmt.Errorf("response parse error")
	}
	if len(*halls) == 0 {
		return halls, fmt.Errorf("hall not found")
	}

	return halls, nil
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

func getRows(query map[string]string) (*model.RackRows, error) {
	rows := new(model.RackRows)
	queryString := ""
	for key, val := range query {
		queryString = queryString + "&" + key + "=" + val
	}
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/row?"+queryString, []byte{})
	if err != nil {
		return rows, err
	}
	if err := json.Unmarshal(body, rows); err != nil {
		return rows, fmt.Errorf("response parse error")
	}
	if len(*rows) == 0 {
		return rows, fmt.Errorf("row not found")
	}

	return rows, nil
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

func getRacks(query map[string]string) (*model.Racks, error) {
	racks := new(model.Racks)
	queryString := ""
	for key, val := range query {
		queryString = queryString + "&" + key + "=" + val
	}
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/rack?"+queryString, []byte{})
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

func getUPSs(query map[string]string) (*model.UPSs, error) {
	upss := new(model.UPSs)
	queryString := ""
	for key, val := range query {
		queryString = queryString + "&" + key + "=" + val
	}
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/ups?"+queryString, []byte{})
	if err != nil {
		return &model.UPSs{}, err
	}
	if err := json.Unmarshal(body, upss); err != nil {
		return &model.UPSs{}, fmt.Errorf("response parse error")
	}

	return upss, nil
}

func getRowPDU(id uint) (*model.RowPDU, error) {
	rowPDU := new(model.RowPDU)
	idStr := strconv.Itoa(int(id))
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/row-pdu/"+idStr, []byte{})
	if err != nil {
		return rowPDU, err
	}
	if err := json.Unmarshal(body, rowPDU); err != nil {
		return rowPDU, fmt.Errorf("response parse error")
	}

	return rowPDU, nil
}

func getRowPDUs(query map[string]string) (*model.RowPDUs, error) {
	pdus := new(model.RowPDUs)
	queryString := ""
	for key, val := range query {
		queryString = queryString + "&" + key + "=" + val
	}
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/row-pdu?"+queryString, []byte{})
	if err != nil {
		return &model.RowPDUs{}, err
	}
	if err := json.Unmarshal(body, pdus); err != nil {
		return &model.RowPDUs{}, fmt.Errorf("response parse error")
	}

	return pdus, nil
}

func getRackPDUs(query map[string]string) (*model.RackPDUs, error) {
	pdus := new(model.RackPDUs)
	queryString := ""
	for key, val := range query {
		queryString = queryString + "&" + key + "=" + val
	}
	body, err := sendRequest("GET", Conf.APIServer.URL+"/datacenter/rack-pdu?"+queryString, []byte{})
	if err != nil {
		return &model.RackPDUs{}, err
	}
	if err := json.Unmarshal(body, pdus); err != nil {
		return &model.RackPDUs{}, fmt.Errorf("response parse error")
	}

	return pdus, nil
}

func loadRackLocation(rack *model.Rack) {
	row, _ := getRow(rack.RowID)
	hall, _ := getHall(row.HallID)
	floor, _ := getFloor(hall.FloorID)
	dc, _ := getDataCenter(floor.DataCenterID)
	floor.DataCenter = dc
	hall.Floor = floor
	row.Hall = hall
	rack.RackRow = row
}
