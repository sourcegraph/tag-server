package main

import "sourcegraph.com/sqs/pbtypes"

const (
	EvtTypeModified   = "modified"
	EvtTypeReferenced = "referenced"
)

// -------------------------------------------------------------------
// The structures below this line are copied and pasted from the
// Sourcegraph main repo

type EvtUpdate struct {
	Hashes []string `protobuf:"bytes,1,rep,name=Hashes" json:"Hashes,omitempty"`
	Users  []string `protobuf:"bytes,2,rep,name=Users" json:"Users,omitempty"`
	Event  *Evt     `protobuf:"bytes,3,opt,name=Event" json:"Event,omitempty"`
}

type EvtsPostOpts struct {
	Updates []*EvtUpdate `protobuf:"bytes,1,rep,name=Updates" json:"Updates,omitempty"`
}

type Evt struct {
	ID    string             `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Title string             `protobuf:"bytes,2,opt,name=Title,proto3" json:"Title,omitempty"`
	Body  string             `protobuf:"bytes,3,opt,name=Body,proto3" json:"Body,omitempty"`
	URL   string             `protobuf:"bytes,4,opt,name=URL,proto3" json:"URL,omitempty"`
	Type  string             `protobuf:"bytes,5,opt,name=Type,proto3" json:"Type,omitempty"`
	Time  *pbtypes.Timestamp `protobuf:"bytes,6,opt,name=Time" json:"Time,omitempty"`
}
