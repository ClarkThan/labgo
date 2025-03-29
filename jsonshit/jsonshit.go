package jsonshit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/buger/jsonparser"
	"github.com/bytedance/sonic"
	"github.com/segmentio/ksuid"
	"github.com/tidwall/gjson"
	"github.com/valyala/fastjson"

	"github.com/ClarkThan/labgo/utils"
)

// 使用sync.Pool复用缓冲区
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func MarshalWithPool(v interface{}) ([]byte, error) {
	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)
	buf.Reset()

	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var data = []byte(`{
	"intentions": [],
	"interval": 2,
	"messages": [
		{
			"contentType": "photo",
			"id": "S7U2pgRvSVFV5mivtXp0W",
			"interval": 1.5,
			"operatorMsg": [],
			"src": "https://tenant-assets.meiqiausercontent.com/widget/301177/hR10/OnPVhv1KTanJIaTbT3Qt.jpg"
		},
		{
			"contentType": "text",
			"id": "Rpl84k5unmvp_Ldc0mXNJ",
			"interval": 1.5,
			"operatorMsg": [],
			"text": "\u003cp\u003e1.\u003cstrong\u003e💯【开学季】24新春开学折扣咨询\u003c/strong\u003e\u003c/p\u003e\u003cp\u003e2.留学申请咨询\u003c/p\u003e\u003cp\u003e3.申校背景提升咨询\u003c/p\u003e\u003cp\u003e4.国际课程辅导（学科GPA、IB、AP、IG、Alevel、竞赛、OSSD、EPQ等）\u003c/p\u003e\u003cp\u003e5.其他咨询也可以直接告诉我\u003c/p\u003e\u003cp\u003e\u003c/p\u003e\u003cp\u003e\u003c/p\u003e\u003cp\u003e😉点击按钮或直接告诉我咨询内容都OK\u003c/p\u003e"
		},
		{
			"id": "ucqGxT_phPlG1pqvfumID",
			"operatorMsg": [],
			"text": {
				"_immutable": {
					"allowUndo": true,
					"currentContent": {
						"blockMap": {
							"1ltnl": {
								"characterList": [
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									}
								],
								"data": {
									"nodeAttributes": {}
								},
								"depth": 0,
								"key": "1ltnl",
								"text": "2.留学申请咨询",
								"type": "unstyled"
							},
							"4lmvg": {
								"characterList": [
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									}
								],
								"data": {
									"nodeAttributes": {}
								},
								"depth": 0,
								"key": "4lmvg",
								"text": "😉点击按钮或直接告诉我咨询内容都OK",
								"type": "unstyled"
							},
							"5pgb3": {
								"characterList": [
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									},
									{
										"entity": null,
										"style": [
											"BOLD"
										]
									}
								],
								"data": {
									"nodeAttributes": {}
								},
								"depth": 0,
								"key": "5pgb3",
								"text": "1.💯【开学季】24新春开学折扣咨询",
								"type": "unstyled"
							},
							"7v7bk": {
								"characterList": [],
								"data": {
									"nodeAttributes": {}
								},
								"depth": 0,
								"key": "7v7bk",
								"text": "",
								"type": "unstyled"
							},
							"c1vas": {
								"characterList": [
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									}
								],
								"data": {
									"nodeAttributes": {}
								},
								"depth": 0,
								"key": "c1vas",
								"text": "4.国际课程辅导（学科GPA、IB、AP、IG、Alevel、竞赛、OSSD、EPQ等）",
								"type": "unstyled"
							},
							"emjj2": {
								"characterList": [
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									}
								],
								"data": {
									"nodeAttributes": {}
								},
								"depth": 0,
								"key": "emjj2",
								"text": "3.申校背景提升咨询",
								"type": "unstyled"
							},
							"eucaj": {
								"characterList": [
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									},
									{
										"entity": null,
										"style": []
									}
								],
								"data": {
									"nodeAttributes": {}
								},
								"depth": 0,
								"key": "eucaj",
								"text": "5.其他咨询也可以直接告诉我",
								"type": "unstyled"
							},
							"tksd": {
								"characterList": [],
								"data": {
									"nodeAttributes": {}
								},
								"depth": 0,
								"key": "tksd",
								"text": "",
								"type": "unstyled"
							}
						},
						"entityMap": {},
						"selectionAfter": {
							"anchorKey": "4lmvg",
							"anchorOffset": 19,
							"focusKey": "4lmvg",
							"focusOffset": 19,
							"hasFocus": true,
							"isBackward": false
						},
						"selectionBefore": {
							"anchorKey": "5pgb3",
							"anchorOffset": 0,
							"focusKey": "5pgb3",
							"focusOffset": 0,
							"hasFocus": true,
							"isBackward": false
						}
					},
					"decorator": {
						"decorators": [
							{
								"_decorators": []
							},
							{
								"_decorators": [
									{},
									{}
								]
							}
						]
					},
					"directionMap": {
						"1ltnl": "LTR",
						"4lmvg": "LTR",
						"5pgb3": "LTR",
						"7v7bk": "LTR",
						"c1vas": "LTR",
						"emjj2": "LTR",
						"eucaj": "LTR",
						"tksd": "LTR"
					},
					"forceSelection": false,
					"inCompositionMode": false,
					"inlineStyleOverride": null,
					"lastChangeType": "insert-fragment",
					"nativelyRenderedContent": null,
					"redoStack": [],
					"selection": {
						"anchorKey": "4lmvg",
						"anchorOffset": 19,
						"focusKey": "4lmvg",
						"focusOffset": 19,
						"hasFocus": false,
						"isBackward": false
					},
					"treeMap": {
						"1ltnl": [
							{
								"decoratorKey": null,
								"end": 8,
								"leaves": [
									{
										"end": 8,
										"start": 0
									}
								],
								"start": 0
							}
						],
						"4lmvg": [
							{
								"decoratorKey": null,
								"end": 19,
								"leaves": [
									{
										"end": 19,
										"start": 0
									}
								],
								"start": 0
							}
						],
						"5pgb3": [
							{
								"decoratorKey": null,
								"end": 19,
								"leaves": [
									{
										"end": 2,
										"start": 0
									},
									{
										"end": 19,
										"start": 2
									}
								],
								"start": 0
							}
						],
						"7v7bk": [
							{
								"decoratorKey": null,
								"end": 0,
								"leaves": [
									{
										"end": 0,
										"start": 0
									}
								],
								"start": 0
							}
						],
						"c1vas": [
							{
								"decoratorKey": null,
								"end": 44,
								"leaves": [
									{
										"end": 44,
										"start": 0
									}
								],
								"start": 0
							}
						],
						"emjj2": [
							{
								"decoratorKey": null,
								"end": 10,
								"leaves": [
									{
										"end": 10,
										"start": 0
									}
								],
								"start": 0
							}
						],
						"eucaj": [
							{
								"decoratorKey": null,
								"end": 14,
								"leaves": [
									{
										"end": 14,
										"start": 0
									}
								],
								"start": 0
							}
						],
						"tksd": [
							{
								"decoratorKey": null,
								"end": 0,
								"leaves": [
									{
										"end": 0,
										"start": 0
									}
								],
								"start": 0
							}
						]
					},
					"undoStack": [
						{
							"blockMap": {
								"5pgb3": {
									"characterList": [],
									"data": {},
									"depth": 0,
									"key": "5pgb3",
									"text": "",
									"type": "unstyled"
								}
							},
							"entityMap": {},
							"selectionAfter": {
								"anchorKey": "5pgb3",
								"anchorOffset": 0,
								"focusKey": "5pgb3",
								"focusOffset": 0,
								"hasFocus": false,
								"isBackward": false
							},
							"selectionBefore": {
								"anchorKey": "5pgb3",
								"anchorOffset": 0,
								"focusKey": "5pgb3",
								"focusOffset": 0,
								"hasFocus": false,
								"isBackward": false
							}
						}
					]
				},
				"convertOptions": {
					"fontFamilies": [
						{
							"family": "Arial, Helvetica, sans-serif",
							"name": "Araial"
						},
						{
							"family": "Georgia, serif",
							"name": "Georgia"
						},
						{
							"family": "Impact, serif",
							"name": "Impact"
						},
						{
							"family": "\"Courier New\", Courier, monospace",
							"name": "Monospace"
						},
						{
							"family": "tahoma, arial, \"Hiragino Sans GB\", 宋体, sans-serif",
							"name": "Tahoma"
						}
					]
				}
			}
		}
	],
	"nanoId": "dnDr2FtTRmeBKHpS_BjRB",
	"options": [
		{
			"content": "1.【开学季】24新春开学折扣咨询",
			"id": "lGuSJmMDgX4vp1IUHcea2",
			"port": "c310ff59-b54c-4824-b7e1-8afc69787bfa"
		},
		{
			"content": "2.留学申请",
			"id": "67f120a6-747d-4aab-f1fe-4db38a76f20d",
			"port": "b59952a9-dad1-428b-83b5-9d088bc8f8a5"
		},
		{
			"content": "3.背景提升",
			"id": "9225352b-926a-4382-fc8e-5631b5a0d747",
			"port": "f7ab46c2-e0a1-4edc-a7f8-d08278324ea2"
		},
		{
			"content": "4.国际课程辅导（学科GPA、IB、AP、IG、Alevel、竞赛、等等）",
			"id": "ebba2e6a-a280-461b-d5be-876cdac7a344",
			"port": "a93a7cd4-7dfa-4274-921f-909b2e81b50d"
		},
		{
			"content": "其他咨询",
			"id": "2109492d-10bc-47ff-96bb-9999057861e8",
			"port": "9b9493e0-1b4b-4fff-9cea-167a9e923ea6"
		}
	],
	"selected": false,
	"silentAsking": {
		"port": "0e6a656b-ec92-4c94-9eea-3bdb5730c04e"
	},
	"silentAskingDuration": 12,
	"silentAskingSwitch": "open",
	"title": "首问语",
	"type": "message",
	"updatedAt": "2024-02-19T09:14:10.153Z"
}`)

var dat = []byte(`[
	{
		"contentType": "photo",
		"id": "S7U2pgRvSVFV5mivtXp0W",
		"interval": 1.5,
		"operatorMsg": [],
		"src": "https://tenant-assets.meiqiausercontent.com/widget/301177/hR10/OnPVhv1KTanJIaTbT3Qt.jpg"
	},
	{
		"contentType": "text",
		"id": "Rpl84k5unmvp_Ldc0mXNJ",
		"interval": 1.5,
		"operatorMsg": [],
		"text": "\u003cp\u003e1.\u003cstrong\u003e💯【开学季】24新春开学折扣咨询\u003c/strong\u003e\u003c/p\u003e\u003cp\u003e2.留学申请咨询\u003c/p\u003e\u003cp\u003e3.申校背景提升咨询\u003c/p\u003e\u003cp\u003e4.国际课程辅导（学科GPA、IB、AP、IG、Alevel、竞赛、OSSD、EPQ等）\u003c/p\u003e\u003cp\u003e5.其他咨询也可以直接告诉我\u003c/p\u003e\u003cp\u003e\u003c/p\u003e\u003cp\u003e\u003c/p\u003e\u003cp\u003e😉点击按钮或直接告诉我咨询内容都OK\u003c/p\u003e"
	},
	{
		"id": "ucqGxT_phPlG1pqvfumID",
		"operatorMsg": [],
		"text": {
			"_immutable": {
				"allowUndo": true,
				"currentContent": {
					"blockMap": {
						"1ltnl": {
							"characterList": [
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								}
							],
							"data": {
								"nodeAttributes": {}
							},
							"depth": 0,
							"key": "1ltnl",
							"text": "2.留学申请咨询",
							"type": "unstyled"
						},
						"4lmvg": {
							"characterList": [
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								}
							],
							"data": {
								"nodeAttributes": {}
							},
							"depth": 0,
							"key": "4lmvg",
							"text": "😉点击按钮或直接告诉我咨询内容都OK",
							"type": "unstyled"
						},
						"5pgb3": {
							"characterList": [
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								},
								{
									"entity": null,
									"style": [
										"BOLD"
									]
								}
							],
							"data": {
								"nodeAttributes": {}
							},
							"depth": 0,
							"key": "5pgb3",
							"text": "1.💯【开学季】24新春开学折扣咨询",
							"type": "unstyled"
						},
						"7v7bk": {
							"characterList": [],
							"data": {
								"nodeAttributes": {}
							},
							"depth": 0,
							"key": "7v7bk",
							"text": "",
							"type": "unstyled"
						},
						"c1vas": {
							"characterList": [
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								}
							],
							"data": {
								"nodeAttributes": {}
							},
							"depth": 0,
							"key": "c1vas",
							"text": "4.国际课程辅导（学科GPA、IB、AP、IG、Alevel、竞赛、OSSD、EPQ等）",
							"type": "unstyled"
						},
						"emjj2": {
							"characterList": [
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								}
							],
							"data": {
								"nodeAttributes": {}
							},
							"depth": 0,
							"key": "emjj2",
							"text": "3.申校背景提升咨询",
							"type": "unstyled"
						},
						"eucaj": {
							"characterList": [
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								},
								{
									"entity": null,
									"style": []
								}
							],
							"data": {
								"nodeAttributes": {}
							},
							"depth": 0,
							"key": "eucaj",
							"text": "5.其他咨询也可以直接告诉我",
							"type": "unstyled"
						},
						"tksd": {
							"characterList": [],
							"data": {
								"nodeAttributes": {}
							},
							"depth": 0,
							"key": "tksd",
							"text": "",
							"type": "unstyled"
						}
					},
					"entityMap": {},
					"selectionAfter": {
						"anchorKey": "4lmvg",
						"anchorOffset": 19,
						"focusKey": "4lmvg",
						"focusOffset": 19,
						"hasFocus": true,
						"isBackward": false
					},
					"selectionBefore": {
						"anchorKey": "5pgb3",
						"anchorOffset": 0,
						"focusKey": "5pgb3",
						"focusOffset": 0,
						"hasFocus": true,
						"isBackward": false
					}
				},
				"decorator": {
					"decorators": [
						{
							"_decorators": []
						},
						{
							"_decorators": [
								{},
								{}
							]
						}
					]
				},
				"directionMap": {
					"1ltnl": "LTR",
					"4lmvg": "LTR",
					"5pgb3": "LTR",
					"7v7bk": "LTR",
					"c1vas": "LTR",
					"emjj2": "LTR",
					"eucaj": "LTR",
					"tksd": "LTR"
				},
				"forceSelection": false,
				"inCompositionMode": false,
				"inlineStyleOverride": null,
				"lastChangeType": "insert-fragment",
				"nativelyRenderedContent": null,
				"redoStack": [],
				"selection": {
					"anchorKey": "4lmvg",
					"anchorOffset": 19,
					"focusKey": "4lmvg",
					"focusOffset": 19,
					"hasFocus": false,
					"isBackward": false
				},
				"treeMap": {
					"1ltnl": [
						{
							"decoratorKey": null,
							"end": 8,
							"leaves": [
								{
									"end": 8,
									"start": 0
								}
							],
							"start": 0
						}
					],
					"4lmvg": [
						{
							"decoratorKey": null,
							"end": 19,
							"leaves": [
								{
									"end": 19,
									"start": 0
								}
							],
							"start": 0
						}
					],
					"5pgb3": [
						{
							"decoratorKey": null,
							"end": 19,
							"leaves": [
								{
									"end": 2,
									"start": 0
								},
								{
									"end": 19,
									"start": 2
								}
							],
							"start": 0
						}
					],
					"7v7bk": [
						{
							"decoratorKey": null,
							"end": 0,
							"leaves": [
								{
									"end": 0,
									"start": 0
								}
							],
							"start": 0
						}
					],
					"c1vas": [
						{
							"decoratorKey": null,
							"end": 44,
							"leaves": [
								{
									"end": 44,
									"start": 0
								}
							],
							"start": 0
						}
					],
					"emjj2": [
						{
							"decoratorKey": null,
							"end": 10,
							"leaves": [
								{
									"end": 10,
									"start": 0
								}
							],
							"start": 0
						}
					],
					"eucaj": [
						{
							"decoratorKey": null,
							"end": 14,
							"leaves": [
								{
									"end": 14,
									"start": 0
								}
							],
							"start": 0
						}
					],
					"tksd": [
						{
							"decoratorKey": null,
							"end": 0,
							"leaves": [
								{
									"end": 0,
									"start": 0
								}
							],
							"start": 0
						}
					]
				},
				"undoStack": [
					{
						"blockMap": {
							"5pgb3": {
								"characterList": [],
								"data": {},
								"depth": 0,
								"key": "5pgb3",
								"text": "",
								"type": "unstyled"
							}
						},
						"entityMap": {},
						"selectionAfter": {
							"anchorKey": "5pgb3",
							"anchorOffset": 0,
							"focusKey": "5pgb3",
							"focusOffset": 0,
							"hasFocus": false,
							"isBackward": false
						},
						"selectionBefore": {
							"anchorKey": "5pgb3",
							"anchorOffset": 0,
							"focusKey": "5pgb3",
							"focusOffset": 0,
							"hasFocus": false,
							"isBackward": false
						}
					}
				]
			},
			"convertOptions": {
				"fontFamilies": [
					{
						"family": "Arial, Helvetica, sans-serif",
						"name": "Araial"
					},
					{
						"family": "Georgia, serif",
						"name": "Georgia"
					},
					{
						"family": "Impact, serif",
						"name": "Impact"
					},
					{
						"family": "\"Courier New\", Courier, monospace",
						"name": "Monospace"
					},
					{
						"family": "tahoma, arial, \"Hiragino Sans GB\", 宋体, sans-serif",
						"name": "Tahoma"
					}
				]
			}
		}
	}
]`)

type Button struct {
	Name  string `mapstructure:"name" json:"name"`
	Type  string `mapstructure:"type" json:"type,omitempty"`
	Value string `mapstructure:"value" json:"value"`
	// Style string `mapstructure:"style" json:"style,omitempty"`
}

// RandMsg 随机消息
type RandMsg struct {
	Interval float32 `mapstructure:"interval" json:"interval"`
	Text     string  `mapstructure:"text" json:"text"`
}

// Message 消息
type Message struct {
	ContentType    string     `mapstructure:"contentType" json:"contentType"` // text | photo | randomText
	Interval       float32    `mapstructure:"interval" json:"interval,omitempty"`
	Text           string     `mapstructure:"text" json:"text,omitempty"`
	OperatorMsg    []*Button  `mapstructure:"operatorMsg" json:"operatorMsg,omitempty"`
	RandomMessages []*RandMsg `mapstructure:"randomMessages" json:"randomMessages,omitempty"`
	Src            string     `mapstructure:"src" json:"src,omitempty"`
}

type Option struct {
	PortID  string `mapstructure:"port" json:"-"`
	Content string `mapstructure:"content" json:"content"`
	Next    string `mapstructure:"next" json:"next,omitempty"`
}

type Intention struct {
	ConfigID int64    `mapstructure:"configId" json:"-"`
	PortID   string   `mapstructure:"port" json:"-"`
	Intent   string   `mapstructure:"intent" json:"intent"`
	Keywords []string `mapstructure:"keywords" json:"keywords"`
	Next     string   `mapstructure:"next" json:"next,omitempty"`
}

type P struct {
	PortID string `mapstructure:"port" json:"-"`
	Next   string `mapstructure:"next" json:"-"`
}

type MsgNodeData struct {
	ID                   string       `mapstructure:"id" json:"id"`
	Type                 string       `mapstructure:"type" json:"type"`
	Title                string       `mapstructure:"title" json:"title"`
	Messages             []*Message   `mapstructure:"messages" json:"messages,omitempty"`
	Options              []*Option    `mapstructure:"options" json:"options,omitempty"`
	Intentions           []*Intention `mapstructure:"intentions" json:"intentions,omitempty"`
	ContentType          string       `mapstructure:"contentType" json:"contentType,omitempty"`
	Success              P            `mapstructure:"success" json:"-"`
	Fail                 P            `mapstructure:"fail" json:"-"`
	Interval             float32      `mapstructure:"interval" json:"interval,omitempty"`
	SilentAskingDuration float32      `mapstructure:"silentAskingDuration" json:"silentAskingDuration,omitempty"`
	SilentAsking         P            `mapstructure:"silentAsking" json:"silentAsking,omitempty"`
	SilentAskingSwitch   string       `mapstructure:"silentAskingSwitch" json:"silentAskingSwitch,omitempty"`
	SilentAskingNext     string       `mapstructure:"silent_asking_next" json:"silent_asking_next,omitempty"`
	Next                 string       `mapstructure:"next" json:"next,omitempty"`
}

func MainOld() {
	var msgs []*Message
	err := json.Unmarshal(dat, &msgs)
	if err != nil {
		log.Printf("err: %v\n", err)
		return
	}

	fmt.Printf("%+v\n", msgs[2])

	// msgNode := new(MsgNodeData)
	// err := json.Unmarshal(data, msgNode)
	// if err != nil {
	// 	log.Printf("err: %v\n", err)
	// 	return
	// }

	// fmt.Println(msgNode.Messages)
}

type ClueItem struct {
	Text  string `json:"text"`
	MsgID int64  `json:"msg_id"`
}

type ConvClues struct {
	Tel    []ClueItem `json:"tel,omitempty"`
	Email  []ClueItem `json:"email,omitempty"`
	Weixin []ClueItem `json:"weixin,omitempty"`
}

func (c *ConvClues) Count() int {
	return len(c.Tel) + len(c.Email) + len(c.Weixin)
}

func parseConvCluesByString(clues string) *ConvClues {
	var cc ConvClues
	if clues == "" {
		return &cc
	}

	err := json.Unmarshal(utils.String2Bytes(clues), &cc)
	if err != nil {
		return nil
	}

	return &cc
}

func (c *ConvClues) ExtractClueTexts() []string {
	cnt := c.Count()
	if cnt == 0 {
		return nil
	}

	texts := make([]string, 0, cnt)
	for _, it := range c.Tel {
		texts = append(texts, it.Text)
	}
	for _, it := range c.Email {
		texts = append(texts, it.Text)
	}
	for _, it := range c.Weixin {
		texts = append(texts, it.Text)
	}

	return texts
}

func demo1() {
	m := map[string]any{
		"foo": []int{1, 2, 3},
		// "shit": &ConvClues{
		// 	Tel: []ClueItem{{Text: "0283-1234567", MsgID: 123}},
		// },
	}

	cc := parseConvCluesByString(`{"weixin":[{"text":"foo123","msg_id":123}]}`)
	if cc != nil {
		log.Println("over...")
		m["detail"] = *cc
	}

	dat, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(dat))
}

type Inner struct {
	Addr string `json:"addr"`
	Type string `json:"type"`
}

type Outer struct {
	Name string `json:"name"`
	Age  int    `json:"aget"`
	Inner
}

func demo2() {
	o := Outer{
		Inner: Inner{
			Type: "nba",
			Addr: "USA",
		},
	}
	o.Name = "air"
	o.Age = 23
	dat, err := json.Marshal(&o)
	if err != nil {
		return
	}

	fmt.Println(string(dat))
}

func demo3() {
	trackID := ksuid.New().String()
	n := time.Now()
	fmt.Println(trackID)
	kid, err := ksuid.Parse(trackID)
	if err != nil {
		fmt.Println("damn: ", err)
		return
	}

	fmt.Println(kid.Time().Before(n))
	// fmt.Println(n)
}

func demo4() {
	dat := []byte(`{"id":"2625ec49-2ae1-4613-8bc0-567fe4a67d2f","source":"mpush-robot","spec_version":"1.0","time":1726655054,"type":"mbot_delay_fired","data":{"ent_id":347448,"conv_id":6708932490,"agent_id":2058788,"action":"silent_asking"}}`)
	convID1 := fastjson.GetInt(dat, "data", "conv_id")
	fmt.Println(convID1)

	convID2, err := jsonparser.GetInt(dat, "data", "conv_id")
	fmt.Println(convID2, err)

	value := gjson.Get(string(dat), "data.conv_id")
	println(value.String())
}

func demo5() {
	m1 := make(map[string]any)
	s := `{"id":"2625ec49-2ae1-4613-8bc0-567fe4a67d2f","source":"mpush-robot","spec_version":"1.0","time":1726655054,"type":"mbot_delay_fired","data":{"ent_id":347448,"conv_id":6708932490,"agent_id":2058788,"action":"silent_asking"}}`
	if err := sonic.UnmarshalString(s, &m1); err != nil {
		log.Println("sonic 1 fail", err)
		return
	}
	fmt.Println(m1)

	m2 := make(map[string]any)
	if err := sonic.Unmarshal(utils.String2Bytes(s), &m2); err != nil {
		log.Println("sonic 2 fail", err)
		return
	}
	fmt.Println(m2)
}

type BigInt struct {
	Value int64 `json:"value,string"` // 序列化为字符串
}

func bigTest() {
	// 使用示例
	data := `{"value": "9223372036854775807"}`
	var bi BigInt
	json.Unmarshal([]byte(data), &bi)
	fmt.Println(bi.Value)
}

type MyTime struct {
	time.Time
}

func (t *MyTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(""), nil
	}
	return []byte(t.Time.Format(`"2006-01-02 15:04:05"`)), nil
}

func (t *MyTime) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	dataStr := string(data[1:(len(data) - 1)])
	tt, err := time.Parse("2006-01-02 15:04:05", dataStr)
	if err != nil {
		return err
	}

	t.Time = tt
	return nil
}

type Event struct {
	ID        int    `json:"id"` // 字段重命名
	Name      string `json:"name"`
	CreatedAt MyTime `json:"created_at,omitzero"` // 时间格式化为字符串
}

func demo6() {
	e := Event{
		ID:   1,
		Name: "Meeting",
		// CreatedAt: MyTime{Time: time.Now()},
	}

	// 序列化为JSON（带缩进格式化）
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	jsonStr := `{"id":2,"name":"Bob","created_at":"2025-03-24 10:08:57"}`
	var e1 Event
	if err := json.Unmarshal([]byte(jsonStr), &e1); err != nil {
		panic(err)
	}

	fmt.Println(e1.CreatedAt)
}

func demo7() {
	e := Event{
		ID:        1,
		Name:      "Meeting",
		CreatedAt: MyTime{Time: time.Now()},
	}

	encoded, _ := MarshalWithPool(e)
	io.Copy(os.Stderr, bytes.NewBuffer(encoded))

	bs := bytes.Buffer{}
	if err := json.NewEncoder(&bs).Encode(e); err != nil {
		log.Fatalf("failed to encode: %v", err)
	}
	io.Copy(os.Stdout, &bs)
	bs.Reset()

	bs.WriteString(`{"id":2,"name":"Bob","created_at":"2025-03-24 10:08:57"}`)
	var e2 Event
	if err := json.NewDecoder(&bs).Decode(&e2); err != nil {
		log.Fatalf("failed to decode: %v", err)
	}
	fmt.Println(e2.Name)
}

func demo8() {
	// 示例 JSON 数据，包含对象和数组的混合结构
	data := `{
		"users": [
			{"name": "Alice", "age": 30},
			{"name": "Bob", "age": 25}
		],
		"status": "active"
	}`

	// 创建一个新的 JSON 解码器
	decoder := json.NewDecoder(strings.NewReader(data))

	// 解析 JSON 对象
	for {
		// 检查是否还有更多的 token
		if !decoder.More() {
			break
		}

		// 获取下一个 token
		token, err := decoder.Token()
		if err != nil {
			fmt.Println("Error decoding token:", err)
			return
		}

		// 输出 token
		fmt.Printf("Token: %v\n", token)

		// 如果 token 是一个对象，则进一步解析
		if delim, ok := token.(json.Delim); ok {
			if delim == '{' {
				// 解析对象
				var obj map[string]interface{}
				if err := decoder.Decode(&obj); err != nil {
					fmt.Println("Error decoding object xxx:", err)
					return
				}
				fmt.Printf("Decoded object: %+v\n", obj)
			} else if delim == '[' {
				// 解析数组
				var users []map[string]interface{}
				if err := decoder.Decode(&users); err != nil {
					fmt.Println("Error decoding array yyy:", err)
					return
				}
				fmt.Printf("Decoded array: %+v\n", users)
			}
		}
	}
}

func handleDynamicJSON(data []byte) {
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		panic(err)
	}

	// 类型断言处理字段
	if name, ok := result["name"].(string); ok {
		fmt.Println("Name:", name)
	}

	// 更安全的处理方式：json.RawMessage
	var raw struct {
		Metadata json.RawMessage `json:"metadata"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		log.Fatalf("got error: %v", err)
	}

	// 延迟解析metadata字段
	var meta map[string]string
	if err := json.Unmarshal(raw.Metadata, &meta); err != nil {
		log.Fatalf("got error: %v", err)
	}

	log.Println(meta)
}

func demo9() {
	bs := []byte(`{"name":"Alice","metadata":{"name":"Alice","addr":"Chengdu"}}`)
	handleDynamicJSON(bs)
}

type Item struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

// 处理大型JSON
func processJSONFlow(r io.Reader) {
	decoder := json.NewDecoder(r)

	// 读取起始分隔符（如数组的'['）
	token, err := decoder.Token()
	if err != nil {
		panic(err)
	}
	if delim, ok := token.(json.Delim); ok && delim.String() != "[" {
		panic("JSON data does not start with '['")
	}

	for decoder.More() {
		var item Item
		if err := decoder.Decode(&item); err != nil {
			panic(err)
		}
		// 处理每个item...
		fmt.Println(item.Name, item.Age)
	}

	// 读取结束分隔符（如数组的']'）
	token, err = decoder.Token()
	if err != nil {
		panic(err)
	}
	if delim, ok := token.(json.Delim); ok && delim.String() != "]" {
		panic("JSON data does not end with ']'")
	}
}

func demo10() {
	bs := bytes.NewBuffer([]byte(`[{"name":"Clark","age":23},{"name":"John","age":24}]`))
	processJSONFlow(bs)
}
func Main() {
	demo8()
}
