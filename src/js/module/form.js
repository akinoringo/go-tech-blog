'use strict';

document.addEventListener('DOMContentLoaded', ()=>{
  const inputs = document.getElementsByTagName('input');
  const form = document.forms.namedItem('article-form');
  const saveBtn = document.querySelector('.article-form__save');
  const cancelBtn = document.querySelector('.article-form__cancel');
  const previewOpenBtn = document.querySelector('.article-form__open-preview');
  const previewCloseBtn = document.querySelector('.article-form__close-preview');
  const articleFormBody = document.querySelector('.article-form__body');
  const articleFormPreview = document.querySelector('.article-form__preview');
  const articleFormBodyTextArea = document.querySelector('.article-form__input--body');
  const articleFormPreviewTextArea = document.querySelector('.article-form__preview-body-contents');
  const errors = document.querySelector('.article-form__errors');
  const errorTmpl = document.querySelector('.article-form__error-tmpl').firstElementChild;

  const mode = {method: '', url: ''};
  if (window.location.pathname.endsWith('new')){
    mode.method = 'POST';
    mode.url = '/articles';
  } else if (window.location.pathname.endsWith('edit')) {
    mode.method = 'PATCH';
    mode.url = `/articles/${window.location.pathname.split('/')[2]}`;
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
    // form の「本文」に入力された Markdown を HTML に変換してプレビューに埋め込みます。
    articleFormPreviewTextArea.innerHTML = md.render(articleFormBodyTextArea.value);

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

  // CSRFトークンの取得
  const csrfToken = document.getElementsByName('csrf')[0].content;

  // 保存処理の実行イベント
  saveBtn.addEventListener('click', event => {
    event.preventDefault();

    errors.innerHTML = null;

    const fd = new FormData(form);

    let status;

    fetch(`/api${url}`, {
      method: method,
      headers: {'X-CSRF-Token': csrfToken},
      body: fd,
    })
    .then(res => {
      status = res.status;
      return res.json();
    })
    .then(body => {
      console.log(JSON.stringify(body));

      if (status === 200) {
        window.location.href = url;
      }

      if (body.ValidationErrors) {
        showErrors(body.ValidationErrors);
      }
    })
    .catch(err => console.error(err));
  });

  const showErrors = messages => {
    if (Array.isArray(messages) && messages.length != 0) {
      const fragment = document.createDocumentFragment();

      messages.forEach(message => {
        const frag = document.createDocumentFragment();

        frag.appendChild(errorTmpl.cloneNode(true));

        frag.querySelector('.article-form__error').innerHTML = message;

        fragment.appendChild(frag);
      })

      errors.appendChild(fragment);
    }
  }
});

