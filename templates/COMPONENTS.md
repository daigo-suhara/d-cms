# Tailwind コンポーネントクラス一覧

`templates/layout/base.html` の `<style type="text/tailwindcss">` で定義。
`login.html` には同クラスのサブセットを直接定義済み。

---

## ボタン

| クラス | 用途 | 主なスタイル |
|---|---|---|
| `btn-primary` | アイコン付きの主要アクション（リンク・ボタン） | `inline-flex gap-1.5 px-3.5 py-2 bg-blue-600 text-white rounded-lg` |
| `btn-blue` | テキストのみの送信ボタン | `px-4 py-2 bg-blue-600 text-white rounded-lg` |
| `btn-secondary` | キャンセル・補助アクション | `px-4 py-2 border border-slate-300 text-slate-600 rounded-lg` |
| `btn-sm-edit` | テーブル行の「編集」小ボタン | `px-2.5 py-1 text-xs border border-slate-300 text-slate-600 rounded-md` |
| `btn-sm-delete` | テーブル行の「削除」小ボタン | `px-2.5 py-1 text-xs border border-red-200 text-red-500 rounded-md` |

```html
<!-- アイコン付き主要ボタン -->
<a href="..." class="btn-primary">
  <svg .../>
  新規作成
</a>

<!-- 送信ボタン -->
<button type="submit" class="btn-blue">保存する</button>
<button type="submit" class="btn-blue w-full">全幅ボタン</button>

<!-- キャンセル -->
<a href="..." class="btn-secondary">キャンセル</a>

<!-- テーブル操作列 -->
<a href="..." class="btn-sm-edit">編集</a>
<button ... class="btn-sm-delete">削除</button>
```

---

## カード

| クラス | 用途 |
|---|---|
| `card` | カードの外枠 `bg-white border border-slate-200 rounded-xl overflow-hidden` |
| `card-header` | カードのヘッダー行（タイトル・操作ボタンを横並び可） |
| `card-body` | カード本文エリア（`p-5`） |
| `card-footer` | カード下部のアクション行（`bg-slate-50 flex gap-2`） |

```html
<div class="card">
  <div class="card-header">タイトル</div>
  <!-- ヘッダーに追加クラスも使用可 -->
  <div class="card-header flex items-center justify-between">
    <span>タイトル</span>
    <button ...>操作</button>
  </div>
  <div class="card-body space-y-4">
    <!-- 本文 -->
  </div>
  <div class="card-footer">
    <button class="btn-blue">保存</button>
    <a class="btn-secondary">キャンセル</a>
  </div>
</div>
```

---

## テーブル

| クラス | 用途 |
|---|---|
| `th` | `<th>` 共通スタイル（左寄せ、小文字大文字化、薄背景） |
| `td` | `<td>` 共通スタイル（`px-5 py-3.5`） |

```html
<div class="card">
  <table class="w-full">
    <thead>
      <tr class="border-b border-slate-200">
        <th class="th">名前</th>
        <th class="th w-32">操作</th>
      </tr>
    </thead>
    <tbody class="divide-y divide-slate-100">
      <tr class="hover:bg-slate-50 transition-colors">
        <td class="td">値</td>
        <td class="td">...</td>
      </tr>
    </tbody>
  </table>
</div>
```

---

## フォーム

| クラス | 用途 |
|---|---|
| `form-input` | テキスト・パスワード等の `<input>` 共通スタイル |
| `form-label` | `<label>` 共通スタイル（`block text-sm font-medium text-slate-700 mb-1.5`） |

```html
<div>
  <label class="form-label">
    フィールド名 <span class="text-red-500">*</span>
  </label>
  <input type="text" name="..." class="form-input" placeholder="...">
  <p class="mt-1 text-xs text-slate-400">補足テキスト</p>
</div>

<!-- font-mono を追加する場合 -->
<input type="text" name="slug" class="form-input font-mono">
```

---

## フィードバック

| クラス | 用途 |
|---|---|
| `alert-error` | エラーメッセージバナー（赤背景・アイコン付き） |

```html
{{ if .error }}
<div class="alert-error">
  <svg class="w-4 h-4 mt-0.5 shrink-0" .../>
  {{ .error }}
</div>
{{ end }}
```

---

## レイアウト

| クラス | 用途 |
|---|---|
| `page-header` | ページ上部のタイトルとアクションボタンを横並びにする行 |
| `page-title` | `<h1>` の共通スタイル |
| `breadcrumb` | パンくずリスト行 |
| `empty-state` | データ 0 件時の空状態表示エリア |

```html
<!-- ページヘッダー -->
<div class="page-header">
  <h1 class="page-title">ページタイトル</h1>
  <a href="..." class="btn-primary">...</a>
</div>

<!-- パンくず付きヘッダー（items-start で揃える） -->
<div class="flex items-start justify-between mb-5">
  <div>
    <div class="breadcrumb">
      <a href="..." class="hover:text-blue-600 transition-colors">親</a>
      <span>/</span>
      <span>現在</span>
    </div>
    <h1 class="page-title">ページタイトル</h1>
  </div>
  <a href="..." class="btn-primary">...</a>
</div>

<!-- 空状態 -->
<div class="empty-state">
  <svg class="w-10 h-10 mx-auto mb-3 text-slate-300" .../>
  <h3 class="text-sm font-semibold text-slate-600 mb-1">データがありません</h3>
  <p class="text-sm text-slate-400">...</p>
</div>
```
