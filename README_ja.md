# VRML Inline Expander

[English version](README.md)

VRMLファイル内の`Inline`ノードを展開し、参照先ファイルの内容を`Group`ノードとして埋め込むCLIツールです。

## 概要

このツールは、VRML 2.0形式のファイルに記述された`Inline`ノードを自動的に検出し、参照されているファイルの内容を展開して1つのファイルにまとめます。複数のVRMLファイルに分割されたシーンを、単一のファイルに統合したい場合に便利です。

## 使い方

### インストール

```bash
go build -o vrml-inline-expander
```

### 実行方法

```bash
./vrml-inline-expander <入力ファイル> <出力ファイル>
```

**パラメータ:**
- `<入力ファイル>`: 展開対象のVRMLファイルのパス
- `<出力ファイル>`: 展開後のVRMLファイルの保存先パス

### 実行例

次のようなVRMLファイルがあるとします。

**top.wrl（入力ファイル）:**
```vrml
#VRML V2.0 utf8

Inline {
  url "refered.wrl"
}
```

**refered.wrl（参照先ファイル）:**
```vrml
#VRML V2.0 utf8

Shape {
  geometry Sphere {}
}
```

以下のコマンドで展開します。

```bash
./vrml-inline-expander top.wrl merged.wrl
```

**merged.wrl（出力ファイル）:**
```vrml
#VRML V2.0 utf8
Group { children [Shape { geometry Sphere { } } ] }
```

`Inline`ノードが`Group`ノードに置き換えられ、参照先ファイルの内容が埋め込まれました。

## 対応形式

- **VRMLバージョン**: VRML 2.0
- **文字コード**: UTF-8
- **参照形式**: ローカルファイルパスのみ（HTTP等のネットワーク参照は非対応）

## 制約事項

- `Inline`ノードの`url`フィールドは単一の文字列形式のみ対応（リスト形式は非対応）
- 同一ファイルへの重複参照は想定していません
- ファイル間でのDEF名の衝突がないことを前提としています

## システム構成図
[システム構成図](/system.md)

## ライセンス

MIT License
