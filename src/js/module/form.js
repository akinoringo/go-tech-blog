'use strict';

document.addEventListener('DOMContentLoaded', ()=>{
  const inputs = document.getElementsByTagName('input');
  const form = document.forms.namedItem('article-form');
  const saveBtn = document.querySelector('.artcile-form__save');
  const cancelBtn = document.querySelector('.article-form__cancel');
  const previewOpenBtn = document.querySelector('.article-form__open-preview');
  const previewCloseBtn = document.querySelector('.article-form__close-preview');
  const articleFormBody = document.querySelector('.article-form__body');
  const articleFormPreview = document.querySelector('.article-form__preview');
  const articleFormBodyTextArea = document.querySelector('.article-form__input--body');
  const articleFormPreviewTextArea = document.querySelector('.article-form__preview-body-contents');

  const mode = {method: '', url: ''};
  if (window.location.pathname.endsWith('new')){
    mode.method = 'POST';
    mode.url = '/';
  } else if (window.location.pathname.endsWith('edit')) {
    mode.method = 'PATCH';
    mode.url = `/&${window.location.pathname.split('/')[1]}`;
  }
  const {method, url} = mode;

  for (let elm of inputs) {
    elm.addEventListener('keydown', event => {
      if (event.keyCode && event.keyCode === 13) {
        event.preventDefault();

        return false;
      }
    })
  }
  previewOpenBtn.addEventListener('click', event => {
    // form の「本文」に入力された内容をプレビューにコピーします。
    articleFormPreviewTextArea.innerHTML = articleFormBodyTextArea.value;

    // 入力フォームを非表示にします。
    articleFormBody.style.display = 'none';

    // プレビューを表示します。
    articleFormPreview.style.display = 'grid';
  });

  // プレビューを閉じるイベントを設定します。
  previewCloseBtn.addEventListener('click', event => {
    // 入力フォームを表示します。
    articleFormBody.style.display = 'grid';

    // プレビューを非表示にします。
    articleFormPreview.style.display = 'none';
  });

  // 前のページに戻るイベントを設定します。
  cancelBtn.addEventListener('click', event => {
    // <button> 要素クリック時のデフォルトの挙動をキャンセルします。
    event.preventDefault();

    // URL を指定して画面を遷移させます。
    window.location.href = url;
  });
});